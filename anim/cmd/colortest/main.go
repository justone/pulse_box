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
	animations := map[string]anim.NewAnimFunc{
		"random_single":        anim.NewRandomSinglePixel,
		"random_fill":          anim.NewRandomFill,
		"random_fill_all_fast": anim.NewRandomAllColorsFast,
		"strand_test":          anim.NewStrandTest,
		"pulse_all":            anim.NewPulseAll,
		"theater":              anim.NewTheaterCrawl,
		"white":                anim.NewAllWhite,
	}

	height := flag.Int("height", 10, "height of the LED grid")
	width := flag.Int("width", 10, "width of the LED grid")
	animation := flag.String("anim", "random_single", "animation to run")
	anim_list := flag.Bool("anim-list", false, "show list of animations")
	logFile := flag.String("logfile", "hub.log", "filename to log messages to")

	rand.Seed(time.Now().Unix())
	flag.Parse()

	if *anim_list {
		for name, _ := range animations {
			fmt.Println(name)
		}
		return
	}

	if _, ok := animations[*animation]; !ok {
		fmt.Println("animation not found, try -anim-list to see what is available")
		os.Exit(2)
	}

	f, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer f.Close()

	log.SetOutput(f)

	if err := run(*height, *width, animations[*animation]); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func run(height, width int, anim_func anim.NewAnimFunc) error {

	// driver, err := anim.NewScreenDriver(height, width)
	driver, err := anim.NewBoxDriver(height, width, os.Getenv("SERIAL_PORT"))
	if err != nil {
		return err
	}

	animation, err := anim_func()
	if err != nil {
		return err
	}

	// override animation with a composite for now
	animation, _ = anim.NewCompositeAnimation(
		anim.NewAllWhite,
		[]anim.ComponentAnimationArg{
			{1, 1, 3, 9, anim.NewRandomSinglePixel},
			{1 + 9 + 1, 1, 3, 9, anim.NewStrandTest},
			// {1, 1 + 3 + 1, 3, 9, anim.NewPulseAll},
			{1, 1 + 3 + 1, 3, 9, anim.NewRandomAllColorsFast},
			{1 + 9 + 1, 1 + 3 + 1, 3, 9, anim.NewRandomFill},
		})

	driver.Start(animation)

	<-driver.DoneChan()

	return nil
}
