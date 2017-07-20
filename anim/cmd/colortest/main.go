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

func NewRandomFill() (*anim.StatefulAnimation, error) {
	return anim.NewStatefulAnimation(func(req, resp chan *anim.Grid) {
		var ids []int
		idx := 0
		on := true

		for {
			g := <-req

			if len(ids) == 0 {
				ids = rand.Perm(len(g.LEDs))
			}

			if len(g.LEDs) == idx {
				idx = 0
				on = !on
			} else {
				if on {
					g.LEDs[ids[idx]].R = 250
				} else {
					g.LEDs[ids[idx]].R = 0
				}
				idx++
			}

			resp <- g
		}
	})
}

func NewRandomAllColorsFast() (*anim.StatefulAnimation, error) {

	return anim.NewStatelessAnimation(func(g *anim.Grid) *anim.Grid {
		for i := 0; i < len(g.LEDs)/3; i++ {
			g.LEDs[rand.Intn(len(g.LEDs))].R = int32(rand.Intn(250))
			g.LEDs[rand.Intn(len(g.LEDs))].G = int32(rand.Intn(250))
			g.LEDs[rand.Intn(len(g.LEDs))].B = int32(rand.Intn(250))
		}

		return g
	})
}

func NewStrandTest() (*anim.StatefulAnimation, error) {
	return anim.NewStatefulAnimation(func(req, resp chan *anim.Grid) {
		idx := 0
		color := 0

		for {
			g := <-req

			if len(g.LEDs) == idx {
				color++
				if color == 3 {
					color = 0
				}
				for _, l := range g.LEDs {
					l.R = 0
					l.G = 0
					l.B = 0
				}

				idx = 0
			} else {
				switch color {
				case 0:
					g.LEDs[idx].R = 250
				case 1:
					g.LEDs[idx].G = 250
				case 2:
					g.LEDs[idx].B = 250
				}
				idx++
			}

			resp <- g
		}
	})
}

func run(height, width int) error {
	driver, err := anim.NewScreenDriver(height, width)
	if err != nil {
		return err
	}

	animation, err := NewStrandTest()
	if err != nil {
		return err
	}

	driver.Start(animation)

	<-driver.DoneChan()

	return nil
}
