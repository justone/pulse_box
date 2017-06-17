#!/usr/bin/env python
import sys

from bibliopixel import LEDMatrix
from bibliopixel.animation import MatrixChannelTest
from BiblioPixelAnimations.matrix import (MatrixRain)
from bibliopixel.drivers.driver_base import DriverBase, ChannelOrder
from bibliopixel.drivers.visualizer import DriverVisualizer

import serial


SERIAL_SPEED = 1000000


class GlediatorSerialArduino(DriverBase):
    """Driver for "Glediator" protocol over serial"""

    def __init__(self, serial_dev, num=0, width=0, height=0,
                 c_order=ChannelOrder.RGB, gamma=None):
        super(GlediatorSerialArduino, self).__init__(
            num, width, height, c_order, gamma)
        self.serial_dev = serial_dev
        self.serial = serial.Serial(self.serial_dev, SERIAL_SPEED)

    def update(self, data):
        """
        Args:
            data (list): Pixel data in the format [R0, G0, B0, R1, G1, B1, ...]
        """

        print "updating"
        # Glediator sends a first byte with a value of 1
        self.serial.write(1)
        # Then send all the pixel values
        self.serial.write(data)
        # Wait until all data is sent
        self.serial.flush()


def main():
    WIDTH = 21
    HEIGHT = 9

    if len(sys.argv) < 2 or sys.argv[1] == 'LOCAL':
        print('Using DriverVisualizer')
        driver = DriverVisualizer(width=WIDTH, height=HEIGHT, pixelSize=48)
    else:
        serial_dev = sys.argv[1]
        print('Using GlediatorSerialArduino serial_dev=%s' % serial_dev)
        driver = GlediatorSerialArduino(serial_dev, width=WIDTH, height=HEIGHT)

    led = LEDMatrix(driver)

    anim = MatrixRain.MatrixRain(led)

    try:
        anim.run(fps=10)
    except KeyboardInterrupt:
        # Ctrl+C will exit the animation and turn the LEDs offs
        led.all_off()
        led.update()


if __name__ == '__main__':
    main()
