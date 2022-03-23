package hive

import (
	"log"
	"sync"
)

// Piece enums
const (
	QueenBee = iota
	Beetle
	Spider
	Grasshopper
	SoldierAnt
)

// Piece colour determines which player it belongs to
const (
	Black = iota
	White
	MaxPlayers
)

// Piece represents a creature tile
type Piece struct {
	creature int
	colour   int
}

// Game holds information on the game state
type Game struct {
	positions map[Hex]Piece // Positions occupied by a piece
}

// Copy returns a deep copy of a Game
func (g *Game) Copy() Game {
	positions := map[Hex]Piece{}
	for h, piece := range g.positions {
		positions[h] = Piece{
			creature: piece.creature,
			colour:   piece.colour,
		}
	}
	gg := Game{positions: positions}
	return gg
}

// checkSpaceOccupied returns true if a space is occupied by a piece
func (g *Game) checkSpaceOccupied(h Hex) bool {
	_, ok := g.positions[h]
	return ok
}

// ensureConnected checks if the graph is connected to enforce the one hive rule
func (g *Game) ensureConnected() bool {
	// Get an arbitrary starting node (consider an empty graph to be connected)
	var start Hex
	for k := range g.positions {
		start = k
		break
	}

	// Valid neighbours must have pieces
	neigbourFunc := func(hh Hex) []Hex {
		neighbours := hh.GetAdjacent()
		validNeigbours := []Hex{}
		for _, neighbour := range neighbours {
			if _, ok := g.positions[neighbour]; ok {
				validNeigbours = append(validNeigbours, neighbour)
			}
		}
		return validNeigbours
	}

	// Traverse graph to find all nodes connected to the starting node
	nodesByDepth := BFS(start, g, neigbourFunc, 0)
	visitedCount := 0
	for i := 0; i < len(nodesByDepth); i++ {
		visitedCount += len(nodesByDepth[i])
	}

	// Check that all nodes were found
	return visitedCount == len(g.positions)
}

// BFS performs a Breadth First Search from a starting hex.
// neighbourFunc should return valid neighbours for a given hex.
func BFS(start Hex, g *Game, neighbourFunc func(Hex) []Hex, maxDepth int) [][]Hex {
	visited := map[Hex]struct{}{
		start: struct{}{},
	}

	nodesByDepth := [][]Hex{
		[]Hex{start},
	}

	for depth := 1; maxDepth == 0 || depth <= maxDepth; depth++ {
		// Increase depth level for found hexes
		nodesByDepth = append(nodesByDepth, []Hex{})

		for _, h := range nodesByDepth[depth-1] {
			neighbours := neighbourFunc(h)
			for _, neighbour := range neighbours {
				// Check that the neighbour has not already been visited
				if _, ok := visited[neighbour]; ok {
					continue
				}

				// Mark neighbour as visited and add to slice of nodes at this depth
				visited[neighbour] = struct{}{}
				nodesByDepth[depth] = append(nodesByDepth[depth], neighbour)
			}
		}
		// Break if no further nodes were found
		if len(nodesByDepth[depth]) == 0 {
			break
		}
	}

	// Pad return value if a specific depth was specified
	for len(nodesByDepth) <= maxDepth {
		nodesByDepth = append(nodesByDepth, []Hex{})
	}

	return nodesByDepth
}

// GetAllAvailableMoves returns a map of hexes to all possible moves for a given player colour
func GetAllAvailableMoves(g Game, colour int) map[Hex][]Hex {
	moves := map[Hex][]Hex{}
	var mu sync.Mutex
	var wg sync.WaitGroup

	for h, piece := range g.positions {
		if piece.colour == colour {
			hh := h
			wg.Add(1)
			go func() {
				m := GetAvailableMoves(hh, g.Copy())
				if len(m) > 0 {
					mu.Lock()
					moves[hh] = m
					mu.Unlock()
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()

	return moves
}

// GetAvailableMoves returns the available moves for a piece
func GetAvailableMoves(h Hex, g Game) []Hex {
	piece, ok := g.positions[h]
	if !ok {
		log.Fatalf("No piece at coordinate %v", h)
	}

	var moves []Hex

	// Remove piece from starting location to avoid invalid moves after the first
	p := g.positions[h]
	delete(g.positions, h)
	// Add piece back at end of function
	defer func() {
		g.positions[h] = p
	}()

	// Check that moving this piece does not break the one hive rule
	if !g.ensureConnected() {
		return moves
	}

	// Get available moves for the piece based on its creature
	switch piece.creature {

	case QueenBee:
		moves = getAvailableAdjacentMoves(h, g, false)

	case Beetle:
		moves = getAvailableAdjacentMoves(h, g, true)

	case Grasshopper:
		moves = getAvailableJumpMoves(h, g)

	case Spider:
		moves = getAvailableBFSMovesAtDepth(h, g, 3)

	case SoldierAnt:
		moves = getAllAvailableBFSMoves(h, g)

	default:
		panic("Unrecognised creature")

	}

	return moves
}

// getAvailableAdjacentMoves returns the available adjacent tiles one away
func getAvailableAdjacentMoves(h Hex, g Game, allowClimbing bool) []Hex {
	adjacent := h.GetAdjacent()
	l := len(adjacent)

	allowed := make([]bool, l)
	for i := 0; i < l; i++ {
		allowed[i] = true
	}

	// Check if each adjacent position is a valid move
	for i := 0; i < l; i++ {
		destHex := adjacent[i]

		prev := i - 1
		if prev < 0 {
			prev = l - 1
		}
		prevHex := adjacent[prev]

		next := i + 1
		if next >= l {
			next = 0
		}
		nextHex := adjacent[next]

		// Check if hexes are occupied
		destHexOccupied := g.checkSpaceOccupied(destHex)
		prevHexOccupied := g.checkSpaceOccupied(prevHex)
		nextHexOccupied := g.checkSpaceOccupied(nextHex)

		// Forbid moves that take the piece out of contact with the hive
		if !destHexOccupied && !prevHexOccupied && !nextHexOccupied {
			allowed[i] = false
		}

		// Forbid moves that are prohibited due to being unable to slide
		if !destHexOccupied && prevHexOccupied && nextHexOccupied {
			allowed[i] = false
		}

		// Exclude moves that would move in to an occupied space unless creature can climb
		if !allowClimbing && destHexOccupied {
			allowed[i] = false
		}
	}

	moves := []Hex{}
	for i, b := range allowed {
		if b {
			moves = append(moves, adjacent[i])
		}
	}
	return moves
}

// getAvailableJumpMoves returns the tiles reachable by jumping over other pieces
func getAvailableJumpMoves(h Hex, g Game) []Hex {
	adjacent := h.GetAdjacent()

	moves := []Hex{}
	for direction := 0; direction < MaxDirections; direction++ {
		adjHex := adjacent[direction]
		if g.checkSpaceOccupied(adjHex) {
			// Move in direction until an empty space is found
			targetHex := adjHex
			for g.checkSpaceOccupied(targetHex) {
				targetHex = targetHex.Move(direction)
			}
			moves = append(moves, targetHex)
		}
	}

	return moves
}

// getAvailableBFSMovesAtDepth returns the available moves the given number of steps away using a BFS
func getAvailableBFSMovesAtDepth(h Hex, g Game, steps int) []Hex {
	neigbourFunc := func(hh Hex) []Hex {
		return getAvailableAdjacentMoves(hh, g, false)
	}
	nodesByDepth := BFS(h, &g, neigbourFunc, steps)
	return nodesByDepth[steps]
}

// getAllAvailableBFSMoves returns all available moves using a BFS
func getAllAvailableBFSMoves(h Hex, g Game) []Hex {
	neigbourFunc := func(hh Hex) []Hex {
		return getAvailableAdjacentMoves(hh, g, false)
	}

	nodesByDepth := BFS(h, &g, neigbourFunc, 0)
	moves := []Hex{}
	for i := 1; i < len(nodesByDepth); i++ {
		moves = append(moves, nodesByDepth[i]...)
	}
	return moves
}
