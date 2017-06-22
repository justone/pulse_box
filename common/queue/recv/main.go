package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	creds := credentials.NewEnvCredentials()

	region := "us-west-2"
	conf := &aws.Config{
		Region:      &region,
		Credentials: creds,
	}

	sess := session.Must(session.NewSession(conf))

	svc := sqs.New(sess)

	url := "https://sqs.us-west-2.amazonaws.com/149259370426/pulse"
	var wait int64 = 20

	receiveMessageInput := &sqs.ReceiveMessageInput{
		QueueUrl:        &url,
		WaitTimeSeconds: &wait,
	}

	output, err := svc.ReceiveMessage(receiveMessageInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(output)

	if len(output.Messages) > 0 {
		deleteMessageInput := &sqs.DeleteMessageInput{
			QueueUrl:      &url,
			ReceiptHandle: output.Messages[0].ReceiptHandle,
		}

		delOutput, err := svc.DeleteMessage(deleteMessageInput)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(delOutput)
	}

}
