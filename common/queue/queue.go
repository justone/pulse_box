package queue

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

type SQSQueue struct {
	sqs *sqs.SQS
	url string
}

type Queue interface {
	Send(interface{}) error
	Receive() (string, error)
}

type SQSConfig struct {
	QueueUrl, Region string
}

var NoMessagesError error = fmt.Errorf("No messages available")

func NewSQS(conf SQSConfig) (*SQSQueue, error) {

	if len(conf.QueueUrl) == 0 {
		return nil, fmt.Errorf("Must provide QueueUrl in SQSConfig")
	}

	if len(conf.Region) == 0 {
		conf.Region = "us-west-2"
	}

	awsConf := &aws.Config{
		Region:      &conf.Region,
		Credentials: credentials.NewEnvCredentials(),
	}

	return &SQSQueue{
		sqs: sqs.New(session.Must(session.NewSession(awsConf))),
		url: conf.QueueUrl,
	}, nil
}

func (q *SQSQueue) Send(thing interface{}) error {

	var message string
	if s, ok := thing.(string); ok {
		message = s
	} else {
		// try to marshal
		b, err := json.Marshal(thing)
		if err != nil {
			return fmt.Errorf("unable to marshal to json: %s", err)
		}
		message = string(b)
	}

	sendMessageInput := &sqs.SendMessageInput{
		QueueUrl:    &q.url,
		MessageBody: &message,
	}
	logrus.Debugf("input to send message: %s", sendMessageInput)

	output, err := q.sqs.SendMessage(sendMessageInput)
	if err != nil {
		return err
	}

	logrus.Debugf("output from send message: %s", output)

	return nil
}

func (q *SQSQueue) Receive() (string, error) {

	var wait int64 = 20
	receiveMessageInput := &sqs.ReceiveMessageInput{
		QueueUrl:        &q.url,
		WaitTimeSeconds: &wait,
	}

	output, err := q.sqs.ReceiveMessage(receiveMessageInput)
	if err != nil {
		return "", err
	}

	logrus.Debugf("output from receive message: %s", output)

	if len(output.Messages) == 0 {
		return "", NoMessagesError
	}

	deleteMessageInput := &sqs.DeleteMessageInput{
		QueueUrl:      &q.url,
		ReceiptHandle: output.Messages[0].ReceiptHandle,
	}

	delOutput, err := q.sqs.DeleteMessage(deleteMessageInput)
	if err != nil {
		return "", err
	}

	logrus.Debugf("output from delete message: %s", delOutput)

	return *output.Messages[0].Body, nil
}

func (q *SQSQueue) ReceiveChan() chan string {
	results := make(chan string)

	go func(res chan string) {
		for {
			str, err := q.Receive()
			if err != nil && err != NoMessagesError {
				logrus.Warnf("error received while receiving message: %s", err)
				close(res)
				return
			}
			res <- str
		}
	}(results)

	return results
}
