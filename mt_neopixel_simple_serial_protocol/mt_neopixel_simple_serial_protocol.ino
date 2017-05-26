#include <Adafruit_NeoPixel.h>
#ifdef __AVR__
  #include <avr/power.h>
#endif

#define NEOPIXEL_DATA_PIN    6
#define NUM_PIXELS        30 * 5  // 30 led/m * 5 m
// Parameter 1 = number of pixels in strip
// Parameter 2 = Arduino pin number (most are valid)
// Parameter 3 = pixel type flags, add together as needed:
//   NEO_KHZ800  800 KHz bitstream (most NeoPixel products w/WS2812 LEDs)
//   NEO_KHZ400  400 KHz (classic 'v1' (not v2) FLORA pixels, WS2811 drivers)
//   NEO_GRB     Pixels are wired for GRB bitstream (most NeoPixel products)
//   NEO_RGB     Pixels are wired for RGB bitstream (v1 FLORA pixels, not v2)
//   NEO_RGBW    Pixels are wired for RGBW bitstream (NeoPixel RGBW products)
Adafruit_NeoPixel strip = Adafruit_NeoPixel(NUM_PIXELS, NEOPIXEL_DATA_PIN, NEO_GRB + NEO_KHZ800);

// IMPORTANT: To reduce NeoPixel burnout risk, add 1000 uF capacitor across
// pixel power leads, add 300 - 500 Ohm resistor on first pixel's data input
// and minimize distance between Arduino and first pixel.  Avoid connecting
// on a live circuit...if you must, connect GND first.



void setup() {
  // This is for Trinket 5V 16MHz, you can remove these three lines if you are not using a Trinket
  #if defined (__AVR_ATtiny85__)
    if (F_CPU == 16000000) clock_prescale_set(clock_div_1);
  #endif
  // End of trinket special code


  strip.begin();
  strip.show(); // Initialize all pixels to 'off'

  Serial.begin(115200);
  dump_serial_buffer();
}

void loop() {
  read_from_serial_and_update_pixels();
//  for (int i = 0; i < NUM_PIXELS; i++) {
//    strip.setPixelColor(i, strip.Color(0,50,0));
//  }
  // Update NeoPixel strip
  strip.show();
  delay(100);
}

void dump_serial_buffer() {
  // Read (and ignore) serial until empty to clear serial buffer
  while (Serial.available() > 1) {
    Serial.read();
  }

}

void read_from_serial_and_update_pixels() {
  if (Serial.available() < 1) {
    Serial.write("no data\n");
    Serial.flush();
    return;
  }

  Serial.write("hello\n");
  Serial.flush();

  // Read serial until we get all pixel values
  // Each pixel needs 3 bytes, so we keep reading until
  // get 3 bytes of values for each pixel
  byte pixel_index = Serial.read();
  byte red = Serial.read();
  byte green = Serial.read();
  byte blue = Serial.read();
  strip.setPixelColor(pixel_index, strip.Color(red, green, blue));

}

