package anim

import (
	"math/rand"
	"time"
)

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

func NewRandomFill() (*StatefulAnimation, error) {
	return NewStatefulAnimation(func(req, resp chan *Grid) {
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

func NewRandomAllColorsFast() (*StatefulAnimation, error) {

	return NewStatelessAnimation(func(g *Grid) *Grid {
		for i := 0; i < len(g.LEDs)/3; i++ {
			g.LEDs[rand.Intn(len(g.LEDs))].R = int32(rand.Intn(250))
			g.LEDs[rand.Intn(len(g.LEDs))].G = int32(rand.Intn(250))
			g.LEDs[rand.Intn(len(g.LEDs))].B = int32(rand.Intn(250))
		}

		return g
	})
}

func NewStrandTest() (*StatefulAnimation, error) {
	return NewStatefulAnimation(func(req, resp chan *Grid) {
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
