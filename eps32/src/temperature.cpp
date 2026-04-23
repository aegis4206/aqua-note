#include <Arduino.h>
#include <OneWire.h>
#include <DallasTemperature.h>

// 設定腳位
#define ONE_WIRE_BUS 2

OneWire oneWire(ONE_WIRE_BUS);
DallasTemperature sensors(&oneWire);

void setup() {
  Serial.begin(115200);
  Serial.println("DS18B20 溫度計測試開始...");

  sensors.begin();
}

void loop() {
  sensors.requestTemperatures(); 
  
  float tempC = sensors.getTempCByIndex(0);

  if (tempC == DEVICE_DISCONNECTED_C) {
    Serial.println("127錯誤：讀取不到溫度！");
  } else {
    Serial.print("目前溫度: ");
    Serial.print(tempC);
    Serial.println(" °C");
  }

  delay(2000);
}
