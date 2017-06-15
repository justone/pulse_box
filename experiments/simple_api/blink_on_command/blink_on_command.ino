/*
  Simple sketch to blink the LED a certain number of times based on serial input.
*/

const int ledPin = LED_BUILTIN;

void setup() {
  pinMode(ledPin, OUTPUT);
  Serial.begin(9600);
}

void loop(){
  if (Serial.available()) {
    int val = Serial.read();
    if (val > '0') {
      light(val - '0');
    }
  }
  delay(500);
}

void light(int n){
  for (int i = 0; i < n; i++) {
    digitalWrite(ledPin, HIGH);
    delay(100);
    digitalWrite(ledPin, LOW);
    delay(100);
  }
}
