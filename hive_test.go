package hive

import (
	"reflect"
	"testing"

	"github.com/maze-mapper/hive/hexgrid"
)

// testCaseGame holds information for a test case game including available moves
type testCaseGame struct {
	game       Game
	moves      map[int]map[hexgrid.Hex][]hexgrid.Hex
	placements map[int][]hexgrid.Hex
}

// zeroPosMoves returns the slice of possible moves for the piece located at the origin
func (tc *testCaseGame) zeroPosMoves() []hexgrid.Hex {
	zeroPos := hexgrid.New(0, 0, 0)
	if moves, ok := tc.moves[Black][zeroPos]; ok {
		return moves
	}
	if moves, ok := tc.moves[White][zeroPos]; ok {
		return moves
	}
	return []hexgrid.Hex{}
}

// Sample games for testing
var sampleGames map[string]testCaseGame = map[string]testCaseGame{
	//      __
	//   __/ B\__
	//  / B\__/ B\
	//  \__/QB\__/
	//  / B\__/
	//  \__/
	"Game 1": {
		game: Game{
			positions: map[hexgrid.Hex]Piece{
				hexgrid.New(0, 0, 0):  Piece{creature: QueenBee},
				hexgrid.New(-1, 1, 0): Piece{creature: Beetle},
				hexgrid.New(-1, 0, 1): Piece{creature: Beetle},
				hexgrid.New(0, -1, 1): Piece{creature: Beetle},
				hexgrid.New(1, -1, 0): Piece{creature: Beetle},
			},
		},
		// Only use this sample for testing central queen bee, moves are incomplete
		moves: map[int]map[hexgrid.Hex][]hexgrid.Hex{
			Black: map[hexgrid.Hex][]hexgrid.Hex{
				hexgrid.New(0, 0, 0): []hexgrid.Hex{hexgrid.New(1, 0, -1), hexgrid.New(0, 1, -1)},
			},
		},
	},

	// Beetle freedom of movement
	//      __
	//   __/ B\__
	//  / B\__/ B\
	//  \__/*B\__/
	//  / B\__/ B\
	//  \__/  \__/
	"Game 2": {
		game: Game{
			positions: map[hexgrid.Hex]Piece{
				hexgrid.New(0, 0, 0):  Piece{creature: Beetle},
				hexgrid.New(-1, 1, 0): Piece{creature: Beetle},
				hexgrid.New(-1, 0, 1): Piece{creature: Beetle},
				hexgrid.New(0, -1, 1): Piece{creature: Beetle},
				hexgrid.New(1, -1, 0): Piece{creature: Beetle},
				hexgrid.New(1, 0, -1): Piece{creature: Beetle},
			},
		},
		// Only use this sample for testing central beetle, moves are incomplete
		moves: map[int]map[hexgrid.Hex][]hexgrid.Hex{
			Black: map[hexgrid.Hex][]hexgrid.Hex{
				hexgrid.New(0, 0, 0): []hexgrid.Hex{
					hexgrid.New(-1, 1, 0), hexgrid.New(-1, 0, 1), hexgrid.New(0, -1, 1), hexgrid.New(1, -1, 0), hexgrid.New(1, 0, -1),
				},
			},
		},
	},

	// Beetle move example
	//     __
	//  __/*B\
	// / S\__/   __
	// \__/SA\__/ G\
	// /QB\__/QB\__/
	// \__/  \__/
	"Game 3": {
		game: Game{
			positions: map[hexgrid.Hex]Piece{
				hexgrid.New(0, 0, 0):   Piece{creature: Beetle, colour: White},
				hexgrid.New(-1, 1, 0):  Piece{creature: Spider, colour: White},
				hexgrid.New(0, 1, -1):  Piece{creature: SoldierAnt, colour: White},
				hexgrid.New(-1, 2, -1): Piece{creature: QueenBee, colour: White},
				hexgrid.New(1, 1, -2):  Piece{creature: QueenBee, colour: Black},
				hexgrid.New(2, 0, -2):  Piece{creature: Grasshopper, colour: Black},
			},
		},
		moves: map[int]map[hexgrid.Hex][]hexgrid.Hex{
			Black: map[hexgrid.Hex][]hexgrid.Hex{
				hexgrid.New(2, 0, -2): []hexgrid.Hex{hexgrid.New(0, 2, -2)},
			},
			White: map[hexgrid.Hex][]hexgrid.Hex{
				hexgrid.New(0, 0, 0):   []hexgrid.Hex{hexgrid.New(-1, 0, 1), hexgrid.New(-1, 1, 0), hexgrid.New(0, 1, -1), hexgrid.New(1, 0, -1)},
				hexgrid.New(-1, 1, 0):  []hexgrid.Hex{hexgrid.New(1, -1, 0), hexgrid.New(-1, 3, -2)},
				hexgrid.New(-1, 2, -1): []hexgrid.Hex{hexgrid.New(-2, 2, 0), hexgrid.New(0, 2, -2)},
			},
		},
		placements: map[int][]hexgrid.Hex{
			Black: []hexgrid.Hex{
				hexgrid.New(1, 2, -3), hexgrid.New(2, 1, -3), hexgrid.New(3, 0, -3), hexgrid.New(3, -1, -2), hexgrid.New(2, -1, -1),
			},
			White: []hexgrid.Hex{
				hexgrid.New(1, -1, 0), hexgrid.New(0, -1, 1), hexgrid.New(-1, 0, 1), hexgrid.New(-2, 1, 1), hexgrid.New(-2, 2, 0), hexgrid.New(-2, 3, -1), hexgrid.New(-1, 3, -2),
			},
		},
	},

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
	"Game 4": {
		game: Game{
			positions: map[hexgrid.Hex]Piece{
				hexgrid.New(0, 0, 0):   Piece{creature: Grasshopper, colour: White},
				hexgrid.New(0, -1, 1):  Piece{creature: QueenBee, colour: Black},
				hexgrid.New(1, -2, 1):  Piece{creature: SoldierAnt, colour: Black},
				hexgrid.New(2, -2, 0):  Piece{creature: SoldierAnt, colour: White},
				hexgrid.New(2, -1, -1): Piece{creature: Spider, colour: White},
				hexgrid.New(1, 0, -1):  Piece{creature: Grasshopper, colour: Black},
				hexgrid.New(0, 1, -1):  Piece{creature: Beetle, colour: White},
				hexgrid.New(2, 0, -2):  Piece{creature: QueenBee, colour: White},
				hexgrid.New(3, -1, -2): Piece{creature: Beetle, colour: Black},
				hexgrid.New(4, -1, -3): Piece{creature: Spider, colour: Black},
				hexgrid.New(4, 0, -4):  Piece{creature: Spider, colour: White},
				hexgrid.New(3, 1, -4):  Piece{creature: Spider, colour: Black},
			},
		},
		// Incomplete move options
		moves: map[int]map[hexgrid.Hex][]hexgrid.Hex{
			White: map[hexgrid.Hex][]hexgrid.Hex{
				hexgrid.New(0, 0, 0): []hexgrid.Hex{
					hexgrid.New(0, -2, 2), hexgrid.New(0, 2, -2), hexgrid.New(3, 0, -3),
				},
			},
		},
	},

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
	"Game 5": {
		game: Game{
			positions: map[hexgrid.Hex]Piece{
				hexgrid.New(0, 0, 0):   Piece{creature: Spider, colour: Black},
				hexgrid.New(1, 0, -1):  Piece{creature: Spider, colour: White},
				hexgrid.New(2, 0, -2):  Piece{creature: QueenBee, colour: White},
				hexgrid.New(2, 1, -3):  Piece{creature: Beetle, colour: Black},
				hexgrid.New(2, 2, -4):  Piece{creature: SoldierAnt, colour: White},
				hexgrid.New(1, 3, -4):  Piece{creature: Grasshopper, colour: Black},
				hexgrid.New(0, 4, -4):  Piece{creature: Grasshopper, colour: White},
				hexgrid.New(-1, 4, -3): Piece{creature: QueenBee, colour: Black},
				hexgrid.New(-1, 3, -2): Piece{creature: SoldierAnt, colour: Black},
				hexgrid.New(-1, 2, -1): Piece{creature: Beetle, colour: White},
			},
		},
		moves: map[int]map[hexgrid.Hex][]hexgrid.Hex{
			Black: map[hexgrid.Hex][]hexgrid.Hex{
				hexgrid.New(0, 0, 0): []hexgrid.Hex{hexgrid.New(3, -1, -2), hexgrid.New(-2, 2, 0), hexgrid.New(0, 3, -3), hexgrid.New(1, 2, -3)},
			},
			White: map[hexgrid.Hex][]hexgrid.Hex{
				hexgrid.New(-1, 2, -1): []hexgrid.Hex{hexgrid.New(-2, 3, -1), hexgrid.New(-1, 3, -2), hexgrid.New(0, 2, -2)},
			},
		},
	},

	// Soldier Ant move example
	//  __    __
	// /QB\  /QB\
	// \__/  \__/
	// / G\__/ B\
	// \__/ B\__/
	//    \__/SA\
	//       \__/
	"Game 6": {
		game: Game{
			positions: map[hexgrid.Hex]Piece{
				hexgrid.New(0, 0, 0):   Piece{creature: SoldierAnt, colour: Black},
				hexgrid.New(0, -1, 1):  Piece{creature: Beetle, colour: Black},
				hexgrid.New(0, -2, 2):  Piece{creature: QueenBee, colour: White},
				hexgrid.New(-1, 0, 1):  Piece{creature: Beetle, colour: White},
				hexgrid.New(-2, 0, 2):  Piece{creature: Grasshopper, colour: White},
				hexgrid.New(-2, -1, 3): Piece{creature: QueenBee, colour: Black},
			},
		},
		moves: map[int]map[hexgrid.Hex][]hexgrid.Hex{
			Black: map[hexgrid.Hex][]hexgrid.Hex{
				hexgrid.New(0, 0, 0): []hexgrid.Hex{
					hexgrid.New(-1, 1, 0), hexgrid.New(-2, 1, 1), hexgrid.New(-3, 1, 2), hexgrid.New(-3, 0, 3), hexgrid.New(-3, -1, 4), hexgrid.New(-2, -2, 4),
					hexgrid.New(-1, -2, 3), hexgrid.New(0, -3, 3), hexgrid.New(1, -3, 2), hexgrid.New(1, -2, 1), hexgrid.New(1, -1, 0),
				},
				hexgrid.New(-2, -1, 3): []hexgrid.Hex{hexgrid.New(-3, 0, 3), hexgrid.New(-1, -1, 2)},
			},
			White: map[hexgrid.Hex][]hexgrid.Hex{
				hexgrid.New(0, -2, 2): []hexgrid.Hex{hexgrid.New(-1, -1, 2), hexgrid.New(1, -2, 1)},
			},
		},
	},

	// One hive example 1
	//  __
	// /QB\__
	// \__/SA\__
	// /SA\__/ B\__
	// \__/  \__/QB\
	//          \__/
	"Game 7": {
		game: Game{
			positions: map[hexgrid.Hex]Piece{
				hexgrid.New(0, 0, 0):  Piece{creature: SoldierAnt, colour: Black},
				hexgrid.New(-1, 0, 1): Piece{creature: QueenBee, colour: Black},
				hexgrid.New(-1, 1, 0): Piece{creature: SoldierAnt, colour: White},
				hexgrid.New(1, 0, -1): Piece{creature: Beetle, colour: Black},
				hexgrid.New(2, 0, -2): Piece{creature: QueenBee, colour: White},
			},
			// Incomplete move options
		},
	},

	// One hive example 2
	//        __
	//     __/SA\__
	//  __/QB\__/ S\
	// / G\__/  \__/
	// \__/   __/ B\
	//       / B\__/
	//       \__/
	"Game 8": {
		game: Game{
			positions: map[hexgrid.Hex]Piece{

				hexgrid.New(0, 0, 0):   Piece{creature: QueenBee, colour: Black},
				hexgrid.New(-1, 1, 0):  Piece{creature: Grasshopper, colour: Black},
				hexgrid.New(1, -1, 0):  Piece{creature: SoldierAnt, colour: Black},
				hexgrid.New(2, -1, -1): Piece{creature: Spider, colour: White},
				hexgrid.New(2, 0, -2):  Piece{creature: Beetle, colour: White},
				hexgrid.New(1, 1, -2):  Piece{creature: Beetle, colour: Black},
			},
		},
		moves: map[int]map[hexgrid.Hex][]hexgrid.Hex{
			Black: map[hexgrid.Hex][]hexgrid.Hex{
				hexgrid.New(-1, 1, 0): []hexgrid.Hex{hexgrid.New(2, -2, 0)},
				hexgrid.New(1, 1, -2): []hexgrid.Hex{hexgrid.New(1, 0, -1), hexgrid.New(2, 0, -2), hexgrid.New(2, 1, -3)},
			},
			White: map[hexgrid.Hex][]hexgrid.Hex{},
		},
	},

	// Trapped spider gets correct moves
	//      __
	//   __/ B\__
	//  / B\__/ B\
	//  \__/*S\__/
	//  / B\__/ B\
	//  \__/  \__/
	"Game 9": {
		game: Game{
			positions: map[hexgrid.Hex]Piece{
				hexgrid.New(0, 0, 0):  Piece{creature: Spider},
				hexgrid.New(-1, 1, 0): Piece{creature: Beetle},
				hexgrid.New(-1, 0, 1): Piece{creature: Beetle},
				hexgrid.New(0, -1, 1): Piece{creature: Beetle},
				hexgrid.New(1, -1, 0): Piece{creature: Beetle},
				hexgrid.New(1, 0, -1): Piece{creature: Beetle},
			},
		},
		// Incomplete move options
	},
}

// TestGetAvailableMoves performs functional tests for getting the available moves for a piece in a game
func TestGetAvailableMoves(t *testing.T) {
	tests := map[string]testCaseGame{
		"Simple slide move":          sampleGames["Game 1"],
		"Beetle freedom of movement": sampleGames["Game 2"],
		"Beetle":                     sampleGames["Game 3"],
		"Grasshopper":                sampleGames["Game 4"],
		"Spider":                     sampleGames["Game 5"],
		"Soldier Ant":                sampleGames["Game 6"],
		"One hive no split":          sampleGames["Game 7"],
		"One hive always linked":     sampleGames["Game 8"],
		"Trapped spider":             sampleGames["Game 9"],
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := GetAvailableMoves(hexgrid.New(0, 0, 0), tc.game)
			want := tc.zeroPosMoves()
			if !hexSlicesAreEqual(got, want) {
				t.Errorf("Got %v, want %v", got, want)
			}
		})
	}
}

func hexSlicesAreEqual(a, b []hexgrid.Hex) bool {
        am := make(map[hexgrid.Hex]struct{})
        for _, v := range a {
                am[v] = struct{}{}
        }

        bm := make(map[hexgrid.Hex]struct{})
        for _, v := range b {
                bm[v] = struct{}{}
        }
        return reflect.DeepEqual(am, bm)
}

func hexDictsAreEqual(a map[hexgrid.Hex][]hexgrid.Hex, b map[hexgrid.Hex][]hexgrid.Hex) bool {
	if len(a) != len(b) {
		return false
	}
	for k, av := range a {
		if bv, ok := b[k]; ok {
			if !hexSlicesAreEqual(av, bv) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

var colourNames map[int]string = map[int]string{
	Black: "black",
	White: "white",
}

func TestGetAllAvailableMoves(t *testing.T) {
	tests := map[string]testCaseGame{
		"Game 3 moves ": sampleGames["Game 3"],
		"Game 5 moves ": sampleGames["Game 5"],
		"Game 6 moves ": sampleGames["Game 6"],
		"Game 8 moves ": sampleGames["Game 8"],
	}
	for name, tc := range tests {
		for player := 0; player < MaxPlayers; player++ {
			t.Run(name+colourNames[player], func(t *testing.T) {
				got := GetAllAvailableMoves(tc.game, player)
				want := tc.moves[player]
				if !hexDictsAreEqual(got, want) {
					t.Errorf("Got %v, want %v", got, want)
				}
			})
		}
	}
}

func TestGetPlacements(t *testing.T) {
	tests := map[string]testCaseGame{
		"Game 3 moves ": sampleGames["Game 3"],
	}
	for name, tc := range tests {
		for player := 0; player < MaxPlayers; player++ {
			t.Run(name+colourNames[player], func(t *testing.T) {
				got := GetPlacements(tc.game, player)
				want := tc.placements[player]
				if !hexSlicesAreEqual(got, want) {
					t.Errorf("Got %v, want %v", got, want)
				}
			})
		}
	}
}
