package hexgrid

import (
	"reflect"
	"testing"
)

func hexSlicesAreEqual(a, b []Hex) bool {
	am := make(map[Hex]struct{})
	for _, v := range a {
		am[v] = struct{}{}
	}

	bm := make(map[Hex]struct{})
	for _, v := range b {
		bm[v] = struct{}{}
	}
	return reflect.DeepEqual(am, bm)
}

func TestMove(t *testing.T) {
	tests := map[string]struct {
		input, want Hex
		direction   int
	}{
		"Move up":         {input: Hex{0, 0, 0}, want: Hex{0, -1, 1}, direction: Up},
		"Move up right":   {input: Hex{0, 0, 0}, want: Hex{1, -1, 0}, direction: UpRight},
		"Move down right": {input: Hex{0, 0, 0}, want: Hex{1, 0, -1}, direction: DownRight},
		"Move down":       {input: Hex{0, 0, 0}, want: Hex{0, 1, -1}, direction: Down},
		"Move down left":  {input: Hex{0, 0, 0}, want: Hex{-1, 1, 0}, direction: DownLeft},
		"Move up left":    {input: Hex{0, 0, 0}, want: Hex{-1, 0, 1}, direction: UpLeft},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.Move(tc.direction)
			if got != tc.want {
				t.Errorf("Got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestGetAdjacent(t *testing.T) {
	input := Hex{0, 0, 0}
	want := []Hex{
		Hex{1, -1, 0},
		Hex{1, 0, -1},
		Hex{0, 1, -1},
		Hex{0, -1, 1},
		Hex{-1, 1, 0},
		Hex{-1, 0, 1},
	}
	t.Run("adjacent", func(t *testing.T) {
		got := input.GetAdjacent()
		if !hexSlicesAreEqual(got, want) {
			t.Errorf("Got %v, want %v", got, want)
		}
	})
}
