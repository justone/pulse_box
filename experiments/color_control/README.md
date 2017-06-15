# Overview

This is a simple Arduino sketch and Python program to demonstrate communication
over serial from the Raspberry Pi to the Arduino.  Specifically, to tell
a strand of LEDs to turn to a particular color.

# Setup

Steps to reproduce:

## On the Arduino

1. Install the color_on_command sketch onto the Arduino
2. Connect the Arduino to the Pi via a USB cable

## On the Pi

1. Make a Python3 virtual env
2. `pip install -r requirements.txt`
3. `SERIAL=/dev/ttyUSB0 ./send_color.py r` to set to red
