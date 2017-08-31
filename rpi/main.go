package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jacobsa/go-serial/serial"
	"github.com/justone/pulse_box/common/queue"
	"github.com/sirupsen/logrus"
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
	SendSerial([]byte)
}

type RealBoard struct {
	Port io.ReadWriteCloser
}

type FakeBoard struct{}

type Command struct {
	Command string `json:"command"`
	Color   string `json:"color"`
}

func (rb *FakeBoard) RandomLED(color Color) error {
	fmt.Println("turning on", color)
	time.Sleep(200 * time.Millisecond)

	return nil
}

func (rb *RealBoard) RandomLED(color Color) error {
	rn := rand.Intn(189)

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

	rb.SendSerial([]byte{0x1, 0x1, byte(rn), byte(red), byte(green), byte(blue)})
	time.Sleep(200 * time.Millisecond)

	return nil
}

func (rb *FakeBoard) SendSerial(data []byte) {

}

func (rb *RealBoard) SendSerial(data []byte) {
	fmt.Println("Sending: ", hex.EncodeToString(data))
	n, err := rb.Port.Write(data)
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}
	fmt.Println("Wrote", n, "bytes.")
}

func main() {

	debug := flag.Bool("debug", false, "show debug output")
	flag.Parse()

	logrus.SetLevel(logrus.InfoLevel)
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	q, err := queue.NewSQS(queue.SQSConfig{
		QueueUrl: sqsQueueUrl,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	var b Board
	if pn := os.Getenv("SERIAL_PORT"); len(pn) > 0 {
		options := serial.OpenOptions{
			PortName:        pn,
			BaudRate:        256000,
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

		time.Sleep(5 * time.Second)
		b = &RealBoard{port}
	} else {
		fmt.Println("No SERIAL_PORT env var found, not sending to real device")
		b = &FakeBoard{}
	}

	_ = q
	// res := q.ReceiveChan()

	// for m := range res {
	// 	fmt.Println(m)
	// 	var cmd Command
	// 	err := json.Unmarshal([]byte(m), &cmd)
	// 	if err != nil {
	// 		logrus.Infof("error unmarshalling data: %s (data: %s)", err, m)
	// 		continue
	// 	}

	// 	fmt.Println(cmd)
	// 	if cmd.Command == "random_led_pulse" {
	// 		switch cmd.Color {
	// 		case "red":
	// 			b.RandomLED(colorRed)
	// 		case "blue":
	// 			b.RandomLED(colorBlue)
	// 		case "green":
	// 			b.RandomLED(colorGreen)
	// 		case "purple":
	// 			b.RandomLED(colorPurple)
	// 		case "yellow":
	// 			b.RandomLED(colorYellow)
	// 		case "orange":
	// 			b.RandomLED(colorOrange)
	// 		}
	// 	}
	// }
	// for {
	// 	// TODO: replace with reading from sqs queue and then send to board
	// 	b.RandomLED(colorRed)
	// }

	for {
		for i := 0; i < 20; i = i + 1 {
			fmt.Println(i)
			data := []byte{0x1}
			for j := 0; j < 189; j++ {
				data = append(data, byte(i), 0x0, byte(i))
			}

			b.SendSerial(data)
			time.Sleep(40 * time.Millisecond)
		}

		for i := 20; i > 0; i = i - 1 {
			fmt.Println(i)
			data := []byte{0x1}
			for j := 0; j < 189; j++ {
				data = append(data, byte(i), 0x0, byte(i))
			}

			b.SendSerial(data)
			time.Sleep(40 * time.Millisecond)
		}
	}
}
