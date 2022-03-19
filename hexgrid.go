package hive

// Hex represents a coordinate on a hexagonal grid
// +s ____
//   /    \
//  /      \ +q
//  \      /
//   \____/
//  +r
type Hex struct {
	q, r, s int
}

// NewHex returns a new Hex
func NewHex(q, r, s int) Hex {
	if q+r+s != 0 {
		panic("Invalid hexgrid coordinates")
	}
	return Hex{q, r, s}
}

// HexDirectionVectors are the unit vectors to move to an adjacent hex
var HexDirectionVectors = [6]Hex{
	Hex{0, -1, 1},
	Hex{1, -1, 0},
	Hex{1, 0, -1},
	Hex{0, 1, -1},
	Hex{-1, 1, 0},
	Hex{-1, 0, 1},
}

// Friendly names for the direction vectors
// Must index to the correct element in HexDirectionVectors
const (
	Up = iota
	UpRight
	DownRight
	Down
	DownLeft
	UpLeft
	MaxDirections
)

// Move returns the Hex in the given direction
func (h *Hex) Move(direction int) Hex {
	vector := HexDirectionVectors[direction]
	return Hex{
		q: h.q + vector.q,
		r: h.r + vector.r,
		s: h.s + vector.s,
	}
}

// GetAdjacent returns the coordinates of all adjacent hexagons
// Return value is ordered in a clockwise direction
func (h *Hex) GetAdjacent() []Hex {
	adjacent := make([]Hex, MaxDirections)
	for direction := 0; direction < MaxDirections; direction++ {
		adjacent[direction] = h.Move(direction)
	}
	return adjacent
}
