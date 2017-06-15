#!/usr/bin/env python

import serial
import os
import sys
import time

color = sys.argv[1]

serial_dev = os.environ.get('SERIAL', '/dev/ttyUSB0')
# open the serial connection
ser = serial.Serial(serial_dev, int(os.environ.get('SER_RATE', '9600')))
# wait a couple seconds for the arduino to reset
time.sleep(2)

# write color command
ser.write(color)
ser.flush()
