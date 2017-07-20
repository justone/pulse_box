package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

type led struct {
	r, g, b int32
}

type Grid struct {
	leds          []*led
	height, width int
}

type Animation interface {
	RequestChan() chan *Grid
	ResponseChan() chan *Grid
}

type Driver interface {
	Render(Animation)
	DoneChan() chan bool
}

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
	var driver Driver
	var animation Animation

	driver, err := NewScreenDriver(height, width)
	if err != nil {
		return err
	}

	animation, err = NewRandomSinglePixel()
	if err != nil {
		return err
	}

	go driver.Render(animation)

	<-driver.DoneChan()

	return nil
}

type RandomSinglePixel struct {
	request  chan *Grid
	response chan *Grid
}

func NewRandomSinglePixel() (*RandomSinglePixel, error) {

	req := make(chan *Grid)
	res := make(chan *Grid)

	go func(req, resp chan *Grid) {
		color := time.After(time.Duration(rand.Intn(50)) * time.Millisecond)
		var pickLED bool
		for {
			select {
			case g := <-req:
				// log.Println("new frame requested")
				// log.Println(g)
				for _, l := range g.leds {
					if l.r > 0 {
						// log.Printf("  knocking down red on %d", i)
						l.r = l.r - 10
					}
					if l.g > 0 {
						// log.Printf("  knocking down green on %d", i)
						l.g = l.g - 10
					}
					if l.b > 0 {
						// log.Printf("  knocking down blue on %d", i)
						l.b = l.b - 10
					}
				}
				if pickLED {
					switch rand.Intn(3) {
					case 0:
						g.leds[rand.Intn(len(g.leds))].r = 250
					case 1:
						g.leds[rand.Intn(len(g.leds))].b = 250
					case 2:
						g.leds[rand.Intn(len(g.leds))].g = 250
					}
					pickLED = false
				}
				// log.Println("sending frame")
				res <- g
			case <-color:
				pickLED = true
				color = time.After(time.Duration(rand.Intn(50)) * time.Millisecond)
			}
		}
	}(req, res)

	return &RandomSinglePixel{req, res}, nil
}

func (rsp *RandomSinglePixel) RequestChan() chan *Grid {
	return rsp.request
}

func (rsp *RandomSinglePixel) ResponseChan() chan *Grid {
	return rsp.response
}

type ScreenDriver struct {
	screen        tcell.Screen
	done          chan bool
	height, width int
}

func NewScreenDriver(height, width int) (*ScreenDriver, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	err = screen.Init()
	if err != nil {
		return nil, err
	}

	screen.HideCursor()
	blackBase := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlack)
	screen.SetStyle(blackBase)
	screen.Clear()

	return &ScreenDriver{
		screen,
		make(chan bool),
		height,
		width,
	}, nil
}

func (sd *ScreenDriver) DoneChan() chan bool {
	return sd.done
}

func (sd *ScreenDriver) Render(anim Animation) {

	out := anim.RequestChan()
	in := anim.ResponseChan()

	var leds []*led
	log.Println("led count:", sd.height*sd.width)
	for i := 0; i <= sd.height*sd.width; i++ {
		leds = append(leds, &led{})
	}
	grid1 := &Grid{leds, sd.height, sd.width}

	out <- grid1

	eventChan := make(chan tcell.Event)
	go func(screen tcell.Screen, e chan tcell.Event) {
		for {
			event := screen.PollEvent()
			if event == nil {
				return
			}
			log.Println("EVENT1:", event)
			e <- event
		}
	}(sd.screen, eventChan)

	done := make(chan bool)
	go func(e chan tcell.Event, d chan bool) {
		for {
			select {
			case event := <-eventChan:
				log.Printf("EVENT: %T", event)
				switch ev := event.(type) {
				case *tcell.EventKey:
					switch ev.Key() {
					case tcell.KeyCtrlC:
						log.Println("DONE")
						done <- true
						return
					case tcell.KeyRune:
						switch ev.Rune() {
						case 'q':
							log.Println("DONE q")
							done <- true
							return
						}
					}
				}
			}
		}
	}(eventChan, done)

	blackBase := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlack)
	var new *Grid
	for {
		timeout := time.After(50 * time.Millisecond)
		// timeout := time.After(1 * time.Second)
		select {
		case new = <-in:
			// log.Println("received new frame")
			// fmt.Println(new)
		case <-timeout:
			if new != nil {
				// log.Println("showing frame")
				for x := 0; x < new.width; x++ {
					for y := 0; y < new.height; y++ {
						l := new.leds[y*new.width+x]
						c := blackBase.Foreground(tcell.NewRGBColor(l.r, l.g, l.b))
						sd.screen.SetCell(x*2, y, c, '•')
					}
				}
				sd.screen.Show()
			}
			out <- new
			new = nil
		case <-done:
			// clean up screen
			sd.screen.Fini()
			sd.done <- true
			return
		}
	}
}
