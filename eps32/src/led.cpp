
// LED
#include <Arduino.h>
#include <Adafruit_NeoPixel.h>
#include <WiFi.h>

#define PIN 48
#define NUMPIXELS 1
Adafruit_NeoPixel pixels(NUMPIXELS, PIN, NEO_GRB + NEO_KHZ800);

const char *ssid = "white";
const char *password = "00000000";

void setup()
{
  pixels.begin();
  pixels.setPixelColor(0, pixels.Color(0, 255, 0));
  pixels.show();
  Serial.begin(115200);

  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);

  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);

  unsigned long startAttempt = millis();
  while (WiFi.status() != WL_CONNECTED && millis() - startAttempt < 20000)
  {
    delay(500);
    Serial.print(".");
  }
  Serial.println();

  if (WiFi.isConnected())
  {
    Serial.println("WiFi connected!");
    Serial.print("IP address: ");
    Serial.println(WiFi.localIP());
    pixels.setPixelColor(0, pixels.Color(0, 0, 0)); 
    pixels.show();

  }
  else
  {
    Serial.println("WiFi connection failed.");
    pixels.setPixelColor(0, pixels.Color(255, 0, 0)); 
    pixels.show();
    delay(1000);
  }
}

void loop()
{
}