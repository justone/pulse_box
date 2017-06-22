package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

type Color uint8

const (
	colorRed = Color(iota)
	colorGreen
	colorBlue
	colorPurple
	colorYellow
	colorOrange
)

var (
	awsAccessKeyId     = getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey = getenv("AWS_SECRET_ACCESS_KEY")
	sqsQueueUrl        = getenv("SQS_QUEUE_URL")
)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

type Board interface {
	RandomLED(Color) error
}

type RealBoard struct {
	Port io.ReadWriteCloser
}

type FakeBoard struct{}

func (rb *FakeBoard) RandomLED(color Color) error {
	fmt.Println("turning on", color)
	time.Sleep(200 * time.Millisecond)

	return nil
}

func (rb *RealBoard) RandomLED(color Color) error {
	rn := rand.Intn(105)

	// brightness := rand.Intn(150)
	brightness := 200

	var red, green, blue int

	switch color {
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

	rb.sendSerial([]byte{0x1, 0x1, byte(rn), byte(red), byte(green), byte(blue)})
	time.Sleep(200 * time.Millisecond)

	return nil
}

func (rb *RealBoard) sendSerial(data []byte) {
	fmt.Println("Sending: ", hex.EncodeToString(data))
	n, err := rb.Port.Write(data)
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}
	fmt.Println("Wrote", n, "bytes.")
}

func main() {

	var b Board
	if pn := os.Getenv("SERIAL_PORT"); len(pn) > 0 {
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
		b = &RealBoard{port}
	} else {
		fmt.Println("No SERIAL_PORT env var found, not sending to real device")
		b = &FakeBoard{}
	}

	for {
		// TODO: replace with reading from sqs queue and then send to board
		b.RandomLED(colorRed)
	}
}
