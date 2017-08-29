#include <Adafruit_NeoPixel.h>
#ifdef __AVR__
  #include <avr/power.h>
#endif

#define NEOPIXEL_DATA_PIN   6
#define NUM_PIXELS          189
Adafruit_NeoPixel strip = Adafruit_NeoPixel(NUM_PIXELS, NEOPIXEL_DATA_PIN, NEO_GRB + NEO_KHZ800);

#define CMD_NEW_DATA 1

byte cmd[2];
byte single[3];
byte leds[NUM_PIXELS*3];
void setup() {

    // initialize all pixels to off
    strip.begin();

    strip.show();

    for (int i = 0; i < NUM_PIXELS*3; i++) {
      leds[i] = 0;
    }

    Serial.begin(115200);
}

int serialGlediator () {
    while (Serial.available() <= 0) {}
    return Serial.read();
}

int stepDown(int input) {
  if(input == 0 || input < 10) {
    return 0;
  }

  return (int) (input * 6.0 / 8.0);
}

void loop() {
  for (int i = 0; i < NUM_PIXELS; i++) {
    byte red = leds[i*3];
    byte green = leds[i*3+1];
    byte blue = leds[i*3+2];

    strip.setPixelColor(i, strip.Color(red, green, blue));

    leds[i*3] = stepDown(red);
    leds[i*3+1] = stepDown(green);
    leds[i*3+2] = stepDown(blue);
  }

  strip.show();
  delay(100);

  if (!Serial.available()) { return; }

  while (serialGlediator() != CMD_NEW_DATA) {}

  Serial.readBytes((byte*)cmd, 2);
  if (cmd[0] == 1) {
    Serial.readBytes((byte*)single, 3);
    leds[cmd[1]*3] = single[0];
    leds[cmd[1]*3+1] = single[1];
    leds[cmd[1]*3+2] = single[2];
  }
}
