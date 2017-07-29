package anim

import (
	"log"
	"time"

	"github.com/gdamore/tcell"
)

type led struct {
	R, G, B int32
}

type Grid struct {
	LEDs          []*led
	LEDs2D        [][]*led
	height, width int
}

func NewGrid(height, width int) *Grid {
	var LEDs []*led
	LEDs2D := make([][]*led, height)
	for i := range LEDs2D {
		LEDs2D[i] = make([]*led, width)
	}

	log.Println("led count:", height*width)
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			new := &led{}
			LEDs = append(LEDs, new)
			LEDs2D[h][w] = new
		}
	}

	return &Grid{LEDs, LEDs2D, height, width}
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

type StatefulAnimation struct {
	*BaseAnimation
}

func NewStatefulAnimation(f func(chan *Grid, chan *Grid)) (*StatefulAnimation, error) {

	ba := NewBaseAnimation()

	go f(ba.RequestChan(), ba.ResponseChan())

	return &StatefulAnimation{ba}, nil
}

func NewStatelessAnimation(f func(*Grid) *Grid) (*StatefulAnimation, error) {

	return NewStatefulAnimation(func(req, resp chan *Grid) {
		for {
			resp <- f(<-req)
		}
	})
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

	out <- NewGrid(sd.height, sd.width)

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
