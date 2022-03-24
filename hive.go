package hive

import (
	"log"
	"sync"

	"github.com/maze-mapper/hive/hexgrid"
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
	positions map[hexgrid.Hex]Piece // Positions occupied by a piece
}

// Copy returns a deep copy of a Game
func (g *Game) Copy() Game {
	positions := map[hexgrid.Hex]Piece{}
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
func (g *Game) checkSpaceOccupied(h hexgrid.Hex) bool {
	_, ok := g.positions[h]
	return ok
}

// ensureConnected checks if the graph is connected to enforce the one hive rule
func (g *Game) ensureConnected() bool {
	// Get an arbitrary starting node (consider an empty graph to be connected)
	var start hexgrid.Hex
	for k := range g.positions {
		start = k
		break
	}

	// Valid neighbours must have pieces
	neigbourFunc := func(hh hexgrid.Hex) []hexgrid.Hex {
		neighbours := hh.GetAdjacent()
		validNeigbours := []hexgrid.Hex{}
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
func BFS(start hexgrid.Hex, g *Game, neighbourFunc func(hexgrid.Hex) []hexgrid.Hex, maxDepth int) [][]hexgrid.Hex {
	visited := map[hexgrid.Hex]struct{}{
		start: struct{}{},
	}

	nodesByDepth := [][]hexgrid.Hex{
		[]hexgrid.Hex{start},
	}

	for depth := 1; maxDepth == 0 || depth <= maxDepth; depth++ {
		// Increase depth level for found hexes
		nodesByDepth = append(nodesByDepth, []hexgrid.Hex{})

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
		nodesByDepth = append(nodesByDepth, []hexgrid.Hex{})
	}

	return nodesByDepth
}

// GetAllAvailableMoves returns a map of hexes to all possible moves for a given player colour
func GetAllAvailableMoves(g Game, colour int) map[hexgrid.Hex][]hexgrid.Hex {
	moves := map[hexgrid.Hex][]hexgrid.Hex{}
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
func GetAvailableMoves(h hexgrid.Hex, g Game) []hexgrid.Hex {
	piece, ok := g.positions[h]
	if !ok {
		log.Fatalf("No piece at coordinate %v", h)
	}

	var moves []hexgrid.Hex

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
func getAvailableAdjacentMoves(h hexgrid.Hex, g Game, allowClimbing bool) []hexgrid.Hex {
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

	moves := []hexgrid.Hex{}
	for i, b := range allowed {
		if b {
			moves = append(moves, adjacent[i])
		}
	}
	return moves
}

// getAvailableJumpMoves returns the tiles reachable by jumping over other pieces
func getAvailableJumpMoves(h hexgrid.Hex, g Game) []hexgrid.Hex {
	adjacent := h.GetAdjacent()

	moves := []hexgrid.Hex{}
	for direction := 0; direction < hexgrid.MaxDirections; direction++ {
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
func getAvailableBFSMovesAtDepth(h hexgrid.Hex, g Game, steps int) []hexgrid.Hex {
	neigbourFunc := func(hh hexgrid.Hex) []hexgrid.Hex {
		return getAvailableAdjacentMoves(hh, g, false)
	}
	nodesByDepth := BFS(h, &g, neigbourFunc, steps)
	return nodesByDepth[steps]
}

// getAllAvailableBFSMoves returns all available moves using a BFS
func getAllAvailableBFSMoves(h hexgrid.Hex, g Game) []hexgrid.Hex {
	neigbourFunc := func(hh hexgrid.Hex) []hexgrid.Hex {
		return getAvailableAdjacentMoves(hh, g, false)
	}

	nodesByDepth := BFS(h, &g, neigbourFunc, 0)
	moves := []hexgrid.Hex{}
	for i := 1; i < len(nodesByDepth); i++ {
		moves = append(moves, nodesByDepth[i]...)
	}
	return moves
}

// GetPlacements returns all hexes where a particular colour piece could be placed
func GetPlacements(g Game, colour int) []hexgrid.Hex {
	allPlacements := map[hexgrid.Hex]map[int]struct{}{}
	for h, piece := range g.positions {
		neighbours := h.GetAdjacent()
		for _, neighbour := range neighbours {
			// Skip hexes that already contain a piece
			if _, ok := g.positions[neighbour]; ok {
				continue
			}
			// Initialise map
			if _, ok := allPlacements[neighbour]; !ok {
				allPlacements[neighbour] = map[int]struct{}{}
			}
			// Mark what colour pieces this space touches
			allPlacements[neighbour][piece.colour] = struct{}{}
		}
	}

	placements := []hexgrid.Hex{}
	for k, v := range allPlacements {
		// Add to placements if the only touching colour is the player colour
		if _, ok := v[colour]; len(v) == 1 && ok {
			placements = append(placements, k)
		}
	}
	return placements
}
