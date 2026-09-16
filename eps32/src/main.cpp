#include <Arduino.h>
#include <WiFi.h>
#include <PubSubClient.h>
#include <OneWire.h>
#include <DallasTemperature.h>
#include <ArduinoJson.h>
#include <Adafruit_NeoPixel.h>

// 設定腳位
#define ONE_WIRE_BUS 2

#define PIN 48
#define NUMPIXELS 1
Adafruit_NeoPixel pixels(NUMPIXELS, PIN, NEO_GRB + NEO_KHZ800);

OneWire oneWire(ONE_WIRE_BUS);
DallasTemperature sensors(&oneWire);

// WIFI 設定
const char *ssid = "white";
const char *password = "00000000";
// MQTT 設定
const char *mqtt_server = "192.168.0.216";
const int mqtt_port = 1883;
const char *mqtt_user = "iot-test";
const char *mqtt_pass = "iot-test";
const char *mqtt_topic = "aquaponics/sensor/data";
String deviceId;

WiFiClient espClient;
PubSubClient client(espClient);

enum class SysStatus
{
  OK,
  WIFI_FAIL,
  WIFI_CONNECTING,
  MQTT_FAIL,
  PUBLISH_FAIL,
  TEMP_FAIL
};
// 原型宣告
bool setupWifi(unsigned long timeoutMs = 15000);
bool tryConnectMQTT();
void setStatusLED(SysStatus status);
void readTemperatureAndPublish();

void setup()
{
  Serial.begin(115200);
  sensors.begin();
  pixels.begin(); // 初始化 NeoPixel
  pixels.setBrightness(50);

  uint64_t chipid = ESP.getEfuseMac();
  char idBuf[13];
  snprintf(idBuf, sizeof(idBuf), "%04X%08X",
           (uint16_t)(chipid >> 32), (uint32_t)chipid);
  deviceId = String(idBuf);

  if (!setupWifi())
  {
    setStatusLED(SysStatus::WIFI_FAIL);
  }
  client.setServer(mqtt_server, mqtt_port);
}

void loop()
{
  // 檢查 Wi-Fi 連線狀態
  if (WiFi.status() != WL_CONNECTED)
  {
    setStatusLED(SysStatus::WIFI_FAIL);
    if (!setupWifi(5000))
    {
      setStatusLED(SysStatus::WIFI_FAIL);
      delay(5000);
      return;
    }
    setStatusLED(SysStatus::OK);
  }
  // 檢查 MQTT 連線
  if (!client.connected())
  {
    setStatusLED(SysStatus::MQTT_FAIL);
    if (!tryConnectMQTT())
    {
      setStatusLED(SysStatus::MQTT_FAIL);
      delay(5000);
      return;
    }
    setStatusLED(SysStatus::OK);
  }
  // 處理 MQTT 背景任務
  client.loop();

  readTemperatureAndPublish();

  delay(5000);
}

void setStatusLED(SysStatus status)
{
  switch (status)
  {
  case SysStatus::WIFI_FAIL:
    pixels.setPixelColor(0, pixels.Color(255, 0, 0)); // 紅
    break;
  case SysStatus::WIFI_CONNECTING:
    pixels.setPixelColor(0, pixels.Color(0, 255, 0)); // 綠
    break;
  case SysStatus::MQTT_FAIL:
    pixels.setPixelColor(0, pixels.Color(255, 255, 0)); // 黃
    break;
  case SysStatus::PUBLISH_FAIL:
    pixels.setPixelColor(0, pixels.Color(255, 100, 0)); // 橘
    break;
  case SysStatus::TEMP_FAIL:
    pixels.setPixelColor(0, pixels.Color(128, 0, 128)); // 紫
    break;
  case SysStatus::OK:
    pixels.setPixelColor(0, pixels.Color(0, 0, 0)); // 關燈
    break;
  }
  pixels.show();
}

// 連接 Wi-Fi
bool setupWifi(unsigned long timeoutMs)
{
  Serial.print("Connecting to WiFi...");
  WiFi.begin(ssid, password);
  unsigned long startAttemptTime = millis();
  while (WiFi.status() != WL_CONNECTED && millis() - startAttemptTime < timeoutMs)
  {
    setStatusLED(SysStatus::WIFI_CONNECTING);
    setStatusLED(SysStatus::OK);
    delay(500);
    Serial.print(".");
  }

  if (WiFi.status() == WL_CONNECTED)
  {
    Serial.println("\nWiFi connected! IP address: ");
    Serial.println(WiFi.localIP());
    return true;
  }

  Serial.println("\nFailed to connect to WiFi within timeout");
  return false;
}

// 連接 MQTT Broker
bool tryConnectMQTT()
{
  Serial.print("Attempting MQTT connection...");

  if (client.connect(deviceId.c_str(), mqtt_user, mqtt_pass))
  {
    Serial.println("connected");
    return true;
  }
  else
  {
    Serial.print("failed, rc=");
    Serial.println(client.state());
    return false;
  }
}

// 讀取溫度
void readTemperatureAndPublish()
{
  // 讀取溫度
  sensors.requestTemperatures();
  float tempC = sensors.getTempCByIndex(0);

  // 確保讀數有效 (-127 代表感測器錯誤或沒接好)
  if (tempC != DEVICE_DISCONNECTED_C)
  {
    setStatusLED(SysStatus::OK);

    // 建立 JSON 物件
    JsonDocument doc;
    doc["device_code"] = deviceId;
    doc["temperature"] = tempC;

    // 將 JSON 轉為字串
    char jsonBuffer[256];
    serializeJson(doc, jsonBuffer);

    // 發布 MQTT 訊息
    Serial.print("Publishing message: ");
    Serial.println(jsonBuffer);
    if (!client.publish(mqtt_topic, jsonBuffer))
    {
      Serial.println("Error: Failed to publish message");
      setStatusLED(SysStatus::PUBLISH_FAIL);
    }
  }
  else
  {
    Serial.println("Error: Could not read temperature data");
    setStatusLED(SysStatus::TEMP_FAIL);
  }
}