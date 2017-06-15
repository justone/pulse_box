from flask import Flask
import serial
import os
import sys
import time

app = Flask(__name__)

serial_dev = os.environ.get('SERIAL', '/dev/ttyUSB0')
# open the serial connection
ser = serial.Serial(serial_dev, int(os.environ.get('SER_RATE', '9600')))
# wait a couple seconds for the arduino to reset
time.sleep(2)

@app.route('/blink/<character>')
def blink(character):
    ser.write(character.encode('utf8'))
    ser.flush()
    return 'ok'
