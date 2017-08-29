#include <FastLED.h>
#define NUM_LEDS 189
#define DATA_PIN 6
#define CMD_NEW_DATA 1
#define BAUD_RATE 256000

CRGB leds[NUM_LEDS];

void setup() {
  FastLED.addLeds<WS2812B, DATA_PIN, GRB>(leds, NUM_LEDS);
  Serial.begin(BAUD_RATE);
}
int serialGlediator () {
  while (!Serial.available()) {}
  return Serial.read();
}
void loop() {
  while (serialGlediator () != CMD_NEW_DATA) {}
  Serial.readBytes((char*)leds, NUM_LEDS*3);
  FastLED.show();
}
