package hexgrid

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
	Max
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
func (h *Hex) GetAdjacent() []Hex {
	adjacent := make([]Hex, Max)
	for direction := 0; direction < Max; direction++ {
		adjacent[direction] = h.Move(direction)
	}
	return adjacent
}

// BFS performs a Breadth First Search from a staring hex to all hexes a given number of steps away.
// The function f should return true if a neighbouring hex can be moved in to.
func BFS(start Hex, steps int, f func(h Hex) bool) []Hex {
	visited := map[Hex]struct{}{
		start: struct{}{},
	}

	// Set up slice of slices for all hexes n steps away
	fringes := make([][]Hex, steps+1)
	fringes[0] = []Hex{start}
	for i := 1; i <= steps; i++ {
		fringes[i] = []Hex{}
	}

	for step := 1; step <= steps; step++ {
		for _, h := range fringes[step-1] {
			neighbours := h.GetAdjacent()
			for _, neighbour := range neighbours {
				// Check that the neighbour has not already been visited
				if _, ok := visited[neighbour]; ok {
					continue
				}
				// Check that neighbour is a valid move as determined by f
				if ok := f(neighbour); !ok {
					continue
				}
				// Mark neighbour as visited and add to slice of allowed moves
				visited[neighbour] = struct{}{}
				fringes[step] = append(fringes[step], neighbour)
			}
		}
	}
	return fringes[steps]
}
