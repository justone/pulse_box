#!/usr/bin/env python
import random
import sys
import time

import serial

SERIAL_SPEED = 115200
NUM_PIXELS = 30 * 5
MAX = 100


def main():
    serial_dev = sys.argv[1]
    ser = serial.Serial(serial_dev, SERIAL_SPEED)
    last_pixel = None
    while 1:
        print(ser.readline())
        if last_pixel:
            pixel_data = bytes([last_pixel, 0, 0, 0])
            ser.write(pixel_data)
            print('clear', [int(b) for b in pixel_data])
            ser.flush()

        last_pixel = random.randrange(1, NUM_PIXELS)
        pixel_data = bytes([last_pixel,
                            random.randrange(0, MAX),
                            random.randrange(0, MAX),
                            random.randrange(0, MAX)])
        ser.write(pixel_data)
        print('write', [int(b) for b in pixel_data])
        ser.flush()
        time.sleep(1)
    # # Wait until remote says go
    # # ser.readline().strip()

    # for i in range(NUM_PIXELS):
    #     # Write out values for a pixel (3 bytes, RGB)
    #     ser.write(bytes(['a', 'a', 'Z']))


if __name__ == '__main__':
    main()
