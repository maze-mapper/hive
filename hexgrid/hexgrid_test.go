package hexgrid

import (
	"reflect"
	"testing"
)

func hexSliceIsEqual(a, b []Hex) bool {
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
		if !hexSliceIsEqual(got, want) {
			t.Errorf("Got %v, want %v", got, want)
		}
	})
}

func TestBFS(t *testing.T) {
	blocked := map[Hex]struct{}{
		Hex{1, -1, 0}:  struct{}{},
		Hex{2, -1, -1}: struct{}{},
		Hex{2, 0, -2}:  struct{}{},
		Hex{2, 1, -3}:  struct{}{},
		Hex{1, 2, -3}:  struct{}{},
		Hex{1, -1, 0}:  struct{}{},
		Hex{0, 2, -2}:  struct{}{},
		Hex{-1, 2, -1}: struct{}{},
		Hex{-1, 1, 0}:  struct{}{},
		Hex{-2, 1, 1}:  struct{}{},
		Hex{-3, 2, 1}:  struct{}{},
		Hex{-1, -1, 2}: struct{}{},
		Hex{0, -2, 2}:  struct{}{},
		Hex{1, -3, 2}:  struct{}{},
	}
	available := func(h Hex) bool {
		_, ok := blocked[h]
		return !ok
	}

	tests := map[string]struct {
		steps int
		want  []Hex
	}{
		"BFS 1": {steps: 1, want: []Hex{
			Hex{0, -1, 1}, Hex{-1, 0, 1}, Hex{1, 0, -1}, Hex{0, 1, -1},
		}},
		"BFS 2": {steps: 2, want: []Hex{
			Hex{1, -2, 1}, Hex{-2, 0, 2}, Hex{1, 1, -2},
		}},
		"BFS 3": {steps: 3, want: []Hex{
			Hex{2, -3, 1}, Hex{2, -2, 0}, Hex{-3, 1, 2}, Hex{-3, 0, 3}, Hex{-2, -1, 3},
		}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := BFS(Hex{0, 0, 0}, tc.steps, available)
			if !hexSliceIsEqual(got, tc.want) {
				t.Errorf("Got %v, want %v", got, tc.want)
			}
		})
	}
}
