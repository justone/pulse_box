package anim

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestNorm(t *testing.T) {

	log.SetOutput(ioutil.Discard)

	tests := []struct {
		start, steps, max int
		dir               bool
		expected          []int32
	}{
		{0, 5, 100, true, []int32{0, 25, 50, 75, 100}},
		{0, 5, 100, false, []int32{0, 100, 75, 50, 25}},
		{2, 5, 100, true, []int32{50, 75, 100, 0, 25}},
		{0, 12, 100, true, []int32{0, 9, 18, 27, 36, 45, 54, 63, 72, 81, 90, 100}},
	}

	for _, test := range tests {
		rep := &repeater{test.start, test.steps, test.max, test.dir}
		var values []int32
		for i := 0; i < test.steps; i++ {
			values = append(values, rep.Next32())
		}
		if len(test.expected) != len(values) {
			t.Fatalf("Expected len of %d, but got %d", len(test.expected), len(values))
		}
		for i := range test.expected {
			if test.expected[i] != values[i] {
				t.Fatalf("Expected (at index %d) %d, but got %d", i, test.expected[i], values[i])
			}
		}
	}
}
