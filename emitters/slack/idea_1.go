package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
)

func main() {
	token, set := os.LookupEnv("PULSEBOT_SLACK_TOKEN")
	if !set {
		log.Fatalf("Environment variable PULSEBOT_SLACK_TOKEN not set")
	}
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event:")
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)
		default:
			fmt.Printf("Other Event: %v\n", ev)
		}
	}
}
