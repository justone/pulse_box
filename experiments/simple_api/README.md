# Overview

This is a simple Arduino sketch and Python program to demonstrate communication
over serial from the Raspberry Pi to the Arduino.

# Setup

Steps to reproduce:

## On the Arduino

1. Install the blink_on_command sketch onto the Arduino
2. Connect the Arduino to the Pi via a USB cable

## On the Pi

1. Make a Python3 virtual env
2. `pip install -r requirements.txt`
3. `FLASK_APP=app.py SERIAL=/dev/ttyUSB0 flask run`
4. `curl -v localhost:5000/blink/8`
