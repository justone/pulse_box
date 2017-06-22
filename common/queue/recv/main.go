package main

import (
	"flag"
	"fmt"

	"github.com/justone/pulse_box/common/queue"
	"github.com/sirupsen/logrus"
)

func main() {
	queueUrl := flag.String("queue", "", "SQS queue url")
	debug := flag.Bool("debug", false, "show debug output")
	flag.Parse()

	logrus.SetLevel(logrus.InfoLevel)
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	q, err := queue.NewSQS(queue.SQSConfig{
		QueueUrl: *queueUrl,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	str, err := q.Receive()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(str)
}
