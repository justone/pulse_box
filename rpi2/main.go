package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/justone/pulse_box/anim"
	"github.com/justone/pulse_box/common/queue"
	"github.com/sirupsen/logrus"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func run() error {
	height := flag.Int("height", 10, "height of the LED grid")
	width := flag.Int("width", 10, "width of the LED grid")
	logFile := flag.String("logfile", "pulse_box.log", "filename to log messages to")

	f, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	log.SetOutput(f)

	driver, err := anim.NewScreenDriver(*height, *width)
	// driver, err := anim.NewBoxDriver(*height, *width, os.Getenv("SERIAL_PORT"))

	if err != nil {
		return err
	}

	queue, err := queue.NewSQS(queue.SQSConfig{
		QueueUrl: sqsQueueUrl,
	})
	_ = queue

	animation, err := createAnimation(queue)
	// animation, err := anim.NewRandomSinglePixel()

	driver.Start(animation)

	<-driver.DoneChan()

	return nil
}

type Command struct {
	Command string `json:"command"`
	Color   string `json:"color"`
}

func createAnimation(queue *queue.SQSQueue) (anim.Animation, error) {
	queueChan := queue.ReceiveChan()
	return anim.NewStatefulAnimation(func(req, resp chan *anim.Grid) {
		colorsToLight := []string{}
		for {
			select {
			case g := <-req:
				// log.Println("new frame requested")
				// log.Println(g)
				for _, l := range g.LEDs {
					if l.R > 0 {
						// log.Printf("  knocking down red on %d", i)
						l.R = l.R - 5
					}
					if l.G > 0 {
						// log.Printf("  knocking down green on %d", i)
						l.G = l.G - 5
					}
					if l.B > 0 {
						// log.Printf("  knocking down blue on %d", i)
						l.B = l.B - 5
					}
				}
				if len(colorsToLight) > 0 {
					for _, color := range colorsToLight {
						log.Println("Lighting ", color)
						l := g.LEDs[rand.Intn(len(g.LEDs))]
						switch color {
						case "red":
							l.R = 200
							l.G = 0
							l.B = 0
						case "green":
							l.R = 0
							l.G = 200
							l.B = 0
						case "blue":
							l.R = 0
							l.G = 0
							l.B = 200
						case "cyan":
							l.R = 0
							l.G = 200
							l.B = 200
						case "magenta":
							l.R = 200
							l.G = 0
							l.B = 200
						case "yellow":
							l.R = 200
							l.G = 200
							l.B = 0
						}
					}
					colorsToLight = []string{}
				}
				resp <- g
			case m := <-queueChan:
				var cmd Command
				err := json.Unmarshal([]byte(m), &cmd)
				if err != nil {
					logrus.Infof("error unmarshalling data: %s (data: %s)", err, m)
					continue
				}

				log.Println(cmd)
				colorsToLight = append(colorsToLight, cmd.Color)
			}
		}
	})
}

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
