package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/justone/pulse_box/common/queue"
	"github.com/sirupsen/logrus"
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

func main() {
	debug := flag.Bool("debug", false, "show debug output")
	color := flag.String("color", "red", "color to show")
	maxDelay := flag.Int("max", 500, "max delay between messages")
	minDelay := flag.Int("min", 200, "min delay between messages")
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

	rand.Seed(time.Now().Unix())
	for {
		err = q.Send(fmt.Sprintf(`{"command": "random_led_pulse", "color": "%s"}`, *color))
		if err != nil {
			logrus.Warnf("error sending: %s", err)
		}

		delay := rand.Intn(*maxDelay-*minDelay) + *minDelay
		logrus.Infof("sleeping: %d", delay)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}
