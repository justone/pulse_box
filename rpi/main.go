package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

func main() {
	options := serial.OpenOptions{
		PortName:        "/dev/tty.usbserial-A7007bpS",
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	rand.Seed(time.Now().Unix())
	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	time.Sleep(2 * time.Second)

	for {
		randomLED(port)
		// randomLEDFade(port)
		// randomRow(port)
	}

	// b := []byte{}
	// for i := 0; i < 21; i++ {
	// 	b = append(b, 0x42, 0x0, 0x42)
	// }
	// paintRow(port, 0, b)
	// time.Sleep(2 * time.Second)
}

func paintRow(port io.ReadWriteCloser, row int, colors []byte) {

	b := []byte{0x1, 0x3, byte(row)}
	b = append(b, colors...)
	sendSerial(port, b)
}

func sendSerial(port io.ReadWriteCloser, data []byte) {
	fmt.Println("Sending: ", hex.EncodeToString(data))
	n, err := port.Write(data)
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}
	fmt.Println("Wrote", n, "bytes.")
}

func randomLED(port io.ReadWriteCloser) {
	rn := rand.Intn(105)

	// brightness := rand.Intn(150)
	brightness := 200

	var red, green, blue int

	switch color := rand.Intn(6); color {
	case 0:
		red = brightness
		green = 0
		blue = 0
	case 1:
		red = 0
		green = brightness
		blue = 0
	case 2:
		red = 0
		green = 0
		blue = brightness
	case 3:
		red = 0
		green = brightness
		blue = brightness
	case 4:
		red = brightness
		green = 0
		blue = brightness
	case 5:
		red = brightness
		green = brightness
		blue = 0
	}

	sendSerial(port, []byte{0x1, 0x1, byte(rn), byte(red), byte(green), byte(blue)})
	time.Sleep(200 * time.Millisecond)
	// time.Sleep(1 * time.Second)

	// sendSerial(port, []byte{0x1, 0x1, byte(rn), 0x0, 0x0, 0x0})
	// time.Sleep(40 * time.Millisecond)
	// // time.Sleep(1 * time.Second)
}

func randomLEDFade(port io.ReadWriteCloser) {
	rn := rand.Intn(189)
	red := rand.Intn(100)
	green := rand.Intn(100)
	blue := rand.Intn(100)

	for i := 0; i <= 8; i++ {
		sendSerial(port, []byte{
			0x1,
			0x1,
			byte(rn),
			byte((float64(i) / float64(8)) * float64(red)),
			byte((float64(i) / float64(8)) * float64(green)),
			byte((float64(i) / float64(8)) * float64(blue)),
		})
		time.Sleep(40 * time.Millisecond)
		// time.Sleep(1 * time.Second)
	}

	for i := 8; i >= 0; i-- {
		sendSerial(port, []byte{
			0x1,
			0x1,
			byte(rn),
			byte((float64(i) / float64(8)) * float64(red)),
			byte((float64(i) / float64(8)) * float64(green)),
			byte((float64(i) / float64(8)) * float64(blue)),
		})
		time.Sleep(40 * time.Millisecond)
		// time.Sleep(1 * time.Second)
	}

	// sendSerial(port, []byte{0x1, 0x1, byte(rn), 0x0, 0x0, 0x0})
	// time.Sleep(40 * time.Millisecond)
	// // time.Sleep(1 * time.Second)
}

func randomRow(port io.ReadWriteCloser) {
	for i := 0; i < 9; i++ {
		sendSerial(port, []byte{0x1, 0x2, byte(i), 0x42, 0x0, 0x42})
		time.Sleep(80 * time.Millisecond)
		// time.Sleep(1 * time.Second)

		sendSerial(port, []byte{0x1, 0x2, byte(i), 0x0, 0x0, 0x0})
		time.Sleep(40 * time.Millisecond)
		// time.Sleep(1 * time.Second)
	}
}
