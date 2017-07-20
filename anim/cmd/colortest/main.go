package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/justone/pulse_box/anim"
)

func main() {
	height := flag.Int("height", 10, "height of the LED grid")
	width := flag.Int("width", 10, "width of the LED grid")
	logFile := flag.String("logfile", "hub.log", "filename to log messages to")

	rand.Seed(time.Now().Unix())
	flag.Parse()

	f, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer f.Close()

	log.SetOutput(f)

	if err := run(*height, *width); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func run(height, width int) error {
	driver, err := anim.NewScreenDriver(height, width)
	if err != nil {
		return err
	}

	animation, err := anim.NewStrandTest()
	if err != nil {
		return err
	}

	driver.Start(animation)

	<-driver.DoneChan()

	return nil
}
