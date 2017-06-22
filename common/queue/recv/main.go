package main

import (
	"fmt"

	"github.com/justone/pulse_box/common/queue"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	q, err := queue.NewSQS(queue.SQSConfig{
		QueueUrl: "https://sqs.us-west-2.amazonaws.com/149259370426/pulse",
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
