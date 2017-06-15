/*
  Simple sketch to run a string of colors down the strand based on serial input.
*/

#include <Adafruit_NeoPixel.h>
#ifdef __AVR__
  #include <avr/power.h>
#endif

#define PIN 6
#define DELAY 20

Adafruit_NeoPixel strip = Adafruit_NeoPixel(189, PIN, NEO_GRB + NEO_KHZ800);

void setup() {
  Serial.begin(9600);

  strip.begin();
  strip.show(); // Initialize all pixels to 'off'
}

void loop(){
  if (Serial.available()) {
    int val = Serial.read();
    if (val == 'r') {
      colorWipe(strip.Color(50, 0, 0), DELAY); // Red
    } else if (val == 'g') {
      colorWipe(strip.Color(0, 50, 0), DELAY); // Green
    } else if (val == 'b') {
      colorWipe(strip.Color(0, 0, 50), DELAY); // Blue
    }
    colorWipe(strip.Color(0, 0, 0), DELAY); // blank out
  }
  delay(500);
}

// Code from strandtest
// Fill the dots one after the other with a color
void colorWipe(uint32_t c, uint8_t wait) {
  for(uint16_t i=0; i<strip.numPixels(); i++) {
    strip.setPixelColor(i, c);
    strip.show();
    delay(wait);
  }
}

