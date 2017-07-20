package anim

import (
	"log"
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

type led struct {
	R, G, B int32
}

type Grid struct {
	LEDs          []*led
	height, width int
}

type Animation interface {
	RequestChan() chan *Grid
	ResponseChan() chan *Grid
}

type Driver interface {
	Start(Animation)
	DoneChan() chan bool
}

type BaseAnimation struct {
	request  chan *Grid
	response chan *Grid
}

func (ba *BaseAnimation) RequestChan() chan *Grid {
	return ba.request
}

func (ba *BaseAnimation) ResponseChan() chan *Grid {
	return ba.response
}

func NewBaseAnimation() *BaseAnimation {
	return &BaseAnimation{
		make(chan *Grid),
		make(chan *Grid),
	}
}

type RandomSinglePixel struct {
	*BaseAnimation
}

func NewRandomSinglePixel() (*RandomSinglePixel, error) {

	ba := NewBaseAnimation()

	go func(req, resp chan *Grid) {
		color := time.After(time.Duration(rand.Intn(50)) * time.Millisecond)
		var pickLED bool
		for {
			select {
			case g := <-req:
				// log.Println("new frame requested")
				// log.Println(g)
				for _, l := range g.LEDs {
					if l.R > 0 {
						// log.Printf("  knocking down red on %d", i)
						l.R = l.R - 10
					}
					if l.G > 0 {
						// log.Printf("  knocking down green on %d", i)
						l.G = l.G - 10
					}
					if l.B > 0 {
						// log.Printf("  knocking down blue on %d", i)
						l.B = l.B - 10
					}
				}
				if pickLED {
					switch rand.Intn(3) {
					case 0:
						g.LEDs[rand.Intn(len(g.LEDs))].R = 250
					case 1:
						g.LEDs[rand.Intn(len(g.LEDs))].B = 250
					case 2:
						g.LEDs[rand.Intn(len(g.LEDs))].G = 250
					}
					pickLED = false
				}
				// log.Println("sending frame")
				resp <- g
			case <-color:
				pickLED = true
				color = time.After(time.Duration(rand.Intn(50)) * time.Millisecond)
			}
		}
	}(ba.RequestChan(), ba.ResponseChan())

	return &RandomSinglePixel{ba}, nil
}

type StatefulAnimation struct {
	*BaseAnimation
}

func NewStatefulAnimation(f func(chan *Grid, chan *Grid)) (*StatefulAnimation, error) {

	ba := NewBaseAnimation()

	go f(ba.RequestChan(), ba.ResponseChan())

	return &StatefulAnimation{ba}, nil
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

func (sd *ScreenDriver) Start(anim Animation) {

	out := anim.RequestChan()
	in := anim.ResponseChan()

	var LEDs []*led
	log.Println("led count:", sd.height*sd.width)
	for i := 0; i <= sd.height*sd.width; i++ {
		LEDs = append(LEDs, &led{})
	}
	grid1 := &Grid{LEDs, sd.height, sd.width}

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

	go func() {
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
							l := new.LEDs[y*new.width+x]
							c := blackBase.Foreground(tcell.NewRGBColor(l.R, l.G, l.B))
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
	}()
}
