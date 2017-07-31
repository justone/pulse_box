package anim

import (
	"math/rand"
	"time"
)

type NewAnimFunc func() (Animation, error)

type RandomSinglePixel struct {
	*BaseAnimation
}

func NewRandomSinglePixel() (Animation, error) {

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

func NewRandomFill() (Animation, error) {
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

func NewRandomAllColorsFast() (Animation, error) {

	return NewStatelessAnimation(func(g *Grid) *Grid {
		for i := 0; i < len(g.LEDs)/3; i++ {
			g.LEDs[rand.Intn(len(g.LEDs))].R = int32(rand.Intn(250))
			g.LEDs[rand.Intn(len(g.LEDs))].G = int32(rand.Intn(250))
			g.LEDs[rand.Intn(len(g.LEDs))].B = int32(rand.Intn(250))
		}

		return g
	})
}

func NewStrandTest() (Animation, error) {
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

func NewPulseAll() (Animation, error) {
	return NewStatefulAnimation(func(req, resp chan *Grid) {
		steps := 20
		current := 1
		max := 240
		up := true

		for {
			g := <-req

			if (up && current == steps) || (!up && current == 1) {
				up = !up
			}

			for _, l := range g.LEDs {
				l.R = int32(max * current / steps)
				l.G = 0
				l.B = 0
			}

			if up {
				current++
			} else {
				current--
			}

			resp <- g
		}
	})
}

type repeater struct {
	cur, maxSteps, max int
	dir                bool // true mean up
}

func (rep *repeater) Next() int {
	ret := rep.cur

	if rep.dir {
		if rep.cur == rep.maxSteps-1 {
			rep.cur = 0
		} else {
			rep.cur++
		}
	} else {
		if rep.cur == 0 {
			rep.cur = rep.maxSteps - 1
		} else {
			rep.cur--
		}
	}

	return int((float32(ret) / float32(rep.maxSteps-1)) * float32(rep.max))
}

func (rep *repeater) Next32() int32 {
	return int32(rep.Next())
}

func NewTheaterCrawl() (Animation, error) {
	return NewStatefulAnimation(func(req, resp chan *Grid) {
		max := 4
		start := max
		skips := 0

		for {
			g := <-req
			skips++
			if skips == 2 {
				rep := &repeater{start, max, 250, true}

				// go around the edge

				// first, across the top
				for i := range g.LEDs2D[0] {
					g.LEDs2D[0][i].SetAll(rep.Next32())
				}

				// then, down the right side
				for i := 1; i < len(g.LEDs2D)-1; i++ {
					g.LEDs2D[i][len(g.LEDs2D[i])-1].SetAll(rep.Next32())
				}

				// then, the other way across the bottom
				last := len(g.LEDs2D) - 1
				for i := len(g.LEDs2D[last]) - 1; i >= 0; i-- {
					g.LEDs2D[last][i].SetAll(rep.Next32())
				}

				// finally, up the left side
				for i := len(g.LEDs2D) - 2; i > 0; i-- {
					g.LEDs2D[i][0].SetAll(rep.Next32())
				}

				start--
				if start < 0 {
					start = max
				}
				skips = 0
			}
			resp <- g
		}
	})
}
