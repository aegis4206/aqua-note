#include <Arduino.h>
#include <WiFi.h>
#include <PubSubClient.h>
#include <OneWire.h>
#include <DallasTemperature.h>
#include <ArduinoJson.h>

// 設定腳位
#define ONE_WIRE_BUS 2

OneWire oneWire(ONE_WIRE_BUS);
DallasTemperature sensors(&oneWire);

// WIFI 設定
const char* ssid = "white";
const char* password = "00000000";
// MQTT 設定
const char* mqtt_server = "192.168.0.216";  
const int mqtt_port = 1883;               
const char* mqtt_user = ""; 
const char* mqtt_pass = "";
const char* mqtt_topic = "aquaponics/sensor/data";

WiFiClient espClient;
PubSubClient client(espClient);

// 連接 Wi-Fi
void setupWifi() {
  Serial.print("Connecting to WiFi...");
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("\nWiFi connected! IP address: ");
  Serial.println(WiFi.localIP());
}

// 連接 MQTT Broker
void connectMQTT() {
  while (!client.connected()) {
    Serial.print("Attempting MQTT connection...");
    // 建立隨機 Client ID
    String clientId = "ESP32S3Client-";
    clientId += String(random(0xffff), HEX);
    
    // 如果有帳號密碼：client.connect(clientId.c_str(), mqtt_user, mqtt_pass)
    // 如果沒有帳號密碼：client.connect(clientId.c_str())
    if (client.connect(clientId.c_str())) {
      Serial.println("connected");
    } else {
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println("try again in 5 seconds");
      delay(5000);
    }
  }
}

void setup() {
  Serial.begin(115200);
  sensors.begin();
  
  setupWifi();
  client.setServer(mqtt_server, mqtt_port);
}

void loop() {
  // 保持 MQTT 連線
  if (!client.connected()) {
    connectMQTT();
  }
  // 處理 MQTT 背景任務
  client.loop(); 

  // 讀取溫度
  sensors.requestTemperatures(); 
  float tempC = sensors.getTempCByIndex(0);

  // 確保讀數有效 (-127 代表感測器錯誤或沒接好)
  if (tempC != DEVICE_DISCONNECTED_C) {
    // 建立 JSON 物件
    StaticJsonDocument<200> doc;
    doc["device"] = "ESP32-S3";
    doc["temperature"] = tempC;
    // 後續若加入 TDS 感測器，可以直接在這裡加上：
    // doc["tds_ppm"] = your_tds_value;

    // 將 JSON 轉為字串
    char jsonBuffer[256];
    serializeJson(doc, jsonBuffer);

    // 發布 MQTT 訊息
    Serial.print("Publishing message: ");
    Serial.println(jsonBuffer);
    client.publish(mqtt_topic, jsonBuffer);
  } else {
    Serial.println("Error: Could not read temperature data");
  }

  delay(5000); 
}