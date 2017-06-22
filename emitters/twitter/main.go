// heavily inspired by https://github.com/campoy/justforfunc/blob/master/14-twitterbot/main.go
package main

import (
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/justone/pulse_box/common/queue"
	"github.com/sirupsen/logrus"
)

var (
	consumerKey       = getenv("TWITTER_CONSUMER_KEY")
	consumerSecret    = getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = getenv("TWITTER_ACCESS_TOKEN_SECRET")
	queueUrl          = getenv("PULSE_SQS_URL")
)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

func main() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	logrus.SetLevel(logrus.DebugLevel)
	log := &logger{logrus.New()}
	api.SetLogger(log)

	q, err := queue.NewSQS(queue.SQSConfig{QueueUrl: queueUrl})

	if err != nil {
		log.Errorf("Error connecting to queue: %v", err)
		os.Exit(1)
	}

	stream := api.PublicStreamFilter(url.Values{
		"track": []string{"@mediatemple", "#mediatemple", "media temple", "Media Temple", "hackathon"},
		// @mediatemple
		"follow": []string{"684983"},
	})

	defer stream.Stop()

	for v := range stream.C {
		t, ok := v.(anaconda.Tweet)
		if !ok {
			log.Warningf("received unexpected value of type %T", v)
			continue
		}

		log.Infof("Tweet: %s", t.Text)
		err = q.Send(`{"command": "random_led_pulse", "color": "green"}`)
		if err != nil {
			log.Errorf("Error sending to queue: %v", err)
		} else {
			log.Infof(" --> Queued")
		}
	}
}

type logger struct {
	*logrus.Logger
}

func (log *logger) Critical(args ...interface{})                 { log.Error(args...) }
func (log *logger) Criticalf(format string, args ...interface{}) { log.Errorf(format, args...) }
func (log *logger) Notice(args ...interface{})                   { log.Info(args...) }
func (log *logger) Noticef(format string, args ...interface{})   { log.Infof(format, args...) }
