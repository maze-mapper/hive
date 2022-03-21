package hive

import (
	"testing"
)

// Sample games for testing
var (
	//      __
	//   __/ B\__
	//  / B\__/ B\
	//  \__/QB\__/
	//  / B\__/
	//  \__/
	game1 Game = Game{
		positions: map[Hex]Piece{
			NewHex(0, 0, 0):  Piece{creature: QueenBee},
			NewHex(-1, 1, 0): Piece{creature: Beetle},
			NewHex(-1, 0, 1): Piece{creature: Beetle},
			NewHex(0, -1, 1): Piece{creature: Beetle},
			NewHex(1, -1, 0): Piece{creature: Beetle},
		},
	}

	// Beetle freedom of movement
	//      __
	//   __/ B\__
	//  / B\__/ B\
	//  \__/*B\__/
	//  / B\__/ B\
	//  \__/  \__/
	game2 Game = Game{
		positions: map[Hex]Piece{
			NewHex(0, 0, 0):  Piece{creature: Beetle},
			NewHex(-1, 1, 0): Piece{creature: Beetle},
			NewHex(-1, 0, 1): Piece{creature: Beetle},
			NewHex(0, -1, 1): Piece{creature: Beetle},
			NewHex(1, -1, 0): Piece{creature: Beetle},
			NewHex(1, 0, -1): Piece{creature: Beetle},
		},
	}
	game2Moves []Hex = []Hex{
		NewHex(-1, 1, 0), NewHex(-1, 0, 1), NewHex(0, -1, 1), NewHex(1, -1, 0), NewHex(1, 0, -1),
	}

	// Beetle move example
	//     __
	//  __/*B\
	// / S\__/   __
	// \__/SA\__/ G\
	// /QB\__/QB\__/
	// \__/  \__/
	game3 Game = Game{
		positions: map[Hex]Piece{
			NewHex(0, 0, 0):   Piece{creature: Beetle, colour: White},
			NewHex(-1, 1, 0):  Piece{creature: Spider, colour: White},
			NewHex(0, 1, -1):  Piece{creature: SoldierAnt, colour: White},
			NewHex(-1, 2, -1): Piece{creature: QueenBee, colour: White},
			NewHex(1, 1, -2):  Piece{creature: QueenBee, colour: Black},
			NewHex(2, 0, -2):  Piece{creature: Grasshopper, colour: Black},
		},
	}
	game3BlackMoves map[Hex][]Hex = map[Hex][]Hex{
		NewHex(2, 0, -2): []Hex{NewHex(0, 2, -2)},
	}
	game3WhiteMoves map[Hex][]Hex = map[Hex][]Hex{
		NewHex(0, 0, 0):   []Hex{NewHex(-1, 0, 1), NewHex(-1, 1, 0), NewHex(0, 1, -1), NewHex(1, 0, -1)},
		NewHex(-1, 1, 0):  []Hex{NewHex(1, -1, 0), NewHex(-1, 3, -2)},
		NewHex(-1, 2, -1): []Hex{NewHex(-2, 2, 0), NewHex(0, 2, -2)},
	}
	game3Moves []Hex = game3WhiteMoves[NewHex(0, 0, 0)]

	// Grasshopper move example
	//     __
	//  __/SA\__
	// /QB\__/SA\
	// \__/  \__/
	// /*G\__/ S\__
	// \__/ G\__/ B\__
	// / B\__/QB\__/ S\
	// \__/  \__/  \__/
	//           __/ S\
	//          / S\__/
	//          \__/
	game4 Game = Game{
		positions: map[Hex]Piece{
			NewHex(0, 0, 0):   Piece{creature: Grasshopper, colour: White},
			NewHex(0, -1, 1):  Piece{creature: QueenBee, colour: Black},
			NewHex(1, -2, 1):  Piece{creature: SoldierAnt, colour: Black},
			NewHex(2, -2, 0):  Piece{creature: SoldierAnt, colour: White},
			NewHex(2, -1, -1): Piece{creature: Spider, colour: White},
			NewHex(1, 0, -1):  Piece{creature: Grasshopper, colour: Black},
			NewHex(0, 1, -1):  Piece{creature: Beetle, colour: White},
			NewHex(2, 0, -2):  Piece{creature: QueenBee, colour: White},
			NewHex(3, -1, -2): Piece{creature: Beetle, colour: Black},
			NewHex(4, -1, -3): Piece{creature: Spider, colour: Black},
			NewHex(4, 0, -4):  Piece{creature: Spider, colour: White},
			NewHex(3, 1, -4):  Piece{creature: Spider, colour: Black},
		},
	}
	game4Moves []Hex = []Hex{
		NewHex(0, -2, 2), NewHex(0, 2, -2), NewHex(3, 0, -3),
	}

	// Spider move example
	//     __
	//    /*S\__
	//    \__/ S\__
	//  __   \__/QB\
	// / B\     \__/
	// \__/     / B\
	// /SA\     \__/
	// \__/   __/SA\
	// /QB\__/ G\__/
	// \__/ G\__/
	//    \__/
	game5 Game = Game{
		positions: map[Hex]Piece{
			NewHex(0, 0, 0):   Piece{creature: Spider, colour: Black},
			NewHex(1, 0, -1):  Piece{creature: Spider, colour: White},
			NewHex(2, 0, -2):  Piece{creature: QueenBee, colour: White},
			NewHex(2, 1, -3):  Piece{creature: Beetle, colour: Black},
			NewHex(2, 2, -4):  Piece{creature: SoldierAnt, colour: White},
			NewHex(1, 3, -4):  Piece{creature: Grasshopper, colour: Black},
			NewHex(0, 4, -4):  Piece{creature: Grasshopper, colour: White},
			NewHex(-1, 4, -3): Piece{creature: QueenBee, colour: Black},
			NewHex(-1, 3, -2): Piece{creature: SoldierAnt, colour: Black},
			NewHex(-1, 2, -1): Piece{creature: Beetle, colour: White},
		},
	}
	game5BlackMoves map[Hex][]Hex = map[Hex][]Hex{
		NewHex(0, 0, 0): []Hex{NewHex(3, -1, -2), NewHex(-2, 2, 0), NewHex(0, 3, -3), NewHex(1, 2, -3)},
	}
	game5WhiteMoves map[Hex][]Hex = map[Hex][]Hex{
		NewHex(-1, 2, -1): []Hex{NewHex(-2, 3, -1), NewHex(-1, 3, -2), NewHex(0, 2, -2)},
	}
	game5Moves []Hex = game5BlackMoves[NewHex(0, 0, 0)]

	// Soldier Ant move example
	//  __    __
	// /QB\  /QB\
	// \__/  \__/
	// / G\__/ B\
	// \__/ B\__/
	//    \__/SA\
	//       \__/
	game6 Game = Game{
		positions: map[Hex]Piece{
			NewHex(0, 0, 0):   Piece{creature: SoldierAnt, colour: Black},
			NewHex(0, -1, 1):  Piece{creature: Beetle, colour: Black},
			NewHex(0, -2, 2):  Piece{creature: QueenBee, colour: White},
			NewHex(-1, 0, 1):  Piece{creature: Beetle, colour: White},
			NewHex(-2, 0, 2):  Piece{creature: Grasshopper, colour: White},
			NewHex(-2, -1, 3): Piece{creature: QueenBee, colour: Black},
		},
	}
	game6BlackMoves map[Hex][]Hex = map[Hex][]Hex{
		NewHex(0, 0, 0): []Hex{
			NewHex(-1, 1, 0), NewHex(-2, 1, 1), NewHex(-3, 1, 2), NewHex(-3, 0, 3), NewHex(-3, -1, 4), NewHex(-2, -2, 4),
			NewHex(-1, -2, 3), NewHex(0, -3, 3), NewHex(1, -3, 2), NewHex(1, -2, 1), NewHex(1, -1, 0),
		},
		NewHex(-2, -1, 3): []Hex{NewHex(-3, 0, 3), NewHex(-1, -1, 2)},
	}
	game6WhiteMoves map[Hex][]Hex = map[Hex][]Hex{
		NewHex(0, -2, 2): []Hex{NewHex(-1, -1, 2), NewHex(1, -2, 1)},
	}
	game6Moves []Hex = game6BlackMoves[NewHex(0, 0, 0)]

	// One hive example 1
	//  __
	// /QB\__
	// \__/SA\__
	// /SA\__/ B\__
	// \__/  \__/QB\
	//          \__/
	game7 Game = Game{
		positions: map[Hex]Piece{
			NewHex(0, 0, 0):  Piece{creature: SoldierAnt, colour: Black},
			NewHex(-1, 0, 1): Piece{creature: QueenBee, colour: Black},
			NewHex(-1, 1, 0): Piece{creature: SoldierAnt, colour: White},
			NewHex(1, 0, -1): Piece{creature: Beetle, colour: Black},
			NewHex(2, 0, -2): Piece{creature: QueenBee, colour: White},
		},
	}

	// One hive example 2
	//        __
	//     __/SA\__
	//  __/QB\__/ S\
	// / G\__/  \__/
	// \__/   __/ B\
	//       / B\__/
	//       \__/
	game8 Game = Game{
		positions: map[Hex]Piece{

			NewHex(0, 0, 0):   Piece{creature: QueenBee, colour: Black},
			NewHex(-1, 1, 0):  Piece{creature: Grasshopper, colour: Black},
			NewHex(1, -1, 0):  Piece{creature: SoldierAnt, colour: Black},
			NewHex(2, -1, -1): Piece{creature: Spider, colour: White},
			NewHex(2, 0, -2):  Piece{creature: Beetle, colour: White},
			NewHex(1, 1, -2):  Piece{creature: Beetle, colour: Black},
		},
	}
	game8BlackMoves map[Hex][]Hex = map[Hex][]Hex{
		NewHex(-1, 1, 0): []Hex{NewHex(2, -2, 0)},
		NewHex(1, 1, -2): []Hex{NewHex(1, 0, -1), NewHex(2, 0, -2), NewHex(2, 1, -3)},
	}
	game8WhiteMoves map[Hex][]Hex = map[Hex][]Hex{}

	// Trapped spider gets correct moves
	//      __
	//   __/ B\__
	//  / B\__/ B\
	//  \__/*S\__/
	//  / B\__/ B\
	//  \__/  \__/
	game9 Game = Game{
		positions: map[Hex]Piece{
			NewHex(0, 0, 0):  Piece{creature: Spider},
			NewHex(-1, 1, 0): Piece{creature: Beetle},
			NewHex(-1, 0, 1): Piece{creature: Beetle},
			NewHex(0, -1, 1): Piece{creature: Beetle},
			NewHex(1, -1, 0): Piece{creature: Beetle},
			NewHex(1, 0, -1): Piece{creature: Beetle},
		},
	}
)

// TestGetAvailableMoves performs functional tests for getting the available moves for a piece in a game
func TestGetAvailableMoves(t *testing.T) {
	tests := map[string]struct {
		g    Game
		want []Hex
	}{
		"Simple slide move":          {g: game1, want: []Hex{NewHex(1, 0, -1), NewHex(0, 1, -1)}},
		"Beetle freedom of movement": {g: game2, want: game2Moves},
		"Beetle":                     {g: game3, want: game3Moves},
		"Grasshopper":                {g: game4, want: game4Moves},
		"Spider":                     {g: game5, want: game5Moves},
		"Soldier Ant":                {g: game6, want: game6Moves},
		"One hive no split":          {g: game7, want: []Hex{}},
		"One hive always linked":     {g: game8, want: []Hex{}},
		"Trapped spider":             {g: game9, want: []Hex{}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := GetAvailableMoves(NewHex(0, 0, 0), tc.g)
			if !HexSliceIsEqual(got, tc.want) {
				t.Errorf("Got %v, want %v", got, tc.want)
			}
		})
	}
}

func hexDictsAreEqual(a map[Hex][]Hex, b map[Hex][]Hex) bool {
	if len(a) != len(b) {
		return false
	}
	for k, av := range a {
		if bv, ok := b[k]; ok {
			if !HexSliceIsEqual(av, bv) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func TestGetAllAvailableMoves(t *testing.T) {
	tests := map[string]struct {
		g      Game
		colour int
		want   map[Hex][]Hex
	}{
		"Game 3 black moves": {g: game3, colour: Black, want: game3BlackMoves},
		"Game 3 white moves": {g: game3, colour: White, want: game3WhiteMoves},
		"Game 5 black moves": {g: game5, colour: Black, want: game5BlackMoves},
		"Game 5 white moves": {g: game5, colour: White, want: game5WhiteMoves},
		"Game 6 black moves": {g: game6, colour: Black, want: game6BlackMoves},
		"Game 6 white moves": {g: game6, colour: White, want: game6WhiteMoves},
		"Game 8 black moves": {g: game8, colour: Black, want: game8BlackMoves},
		"Game 8 white moves": {g: game8, colour: White, want: game8WhiteMoves},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := GetAllAvailableMoves(tc.g, tc.colour)
			if !hexDictsAreEqual(got, tc.want) {
				t.Errorf("Got %v, want %v", got, tc.want)
			}
		})
	}
}
