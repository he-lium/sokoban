package parse_test

import (
	"testing"

	"github.com/he-lium/sokoban/mock"
	"github.com/he-lium/sokoban/parse"
)

func TestBoardParse(t *testing.T) {
	b1, err := mock.BoardMaker3{}.GenBoard()
	if err != nil {
		t.Fatal(err.Error())
	}

	// convert to json
	json, err := parse.BoardToJSON(b1)
	if err != nil {
		t.Fatalf("error writing to JSON: %s", err.Error())
	}

	// parse and convert back to board
	gen := parse.JSONBoard{JSONContent: json}
	b2, err := gen.GenBoard()
	if err != nil {
		t.Fatalf("error parsing JSON: %s", err.Error())
	}

	if b1.Player != b2.Player {
		t.Errorf("b1 player (%d, %d) b2 (%d, %d)", b1.Player.X, b1.Player.Y,
			b2.Player.X, b2.Player.Y)
	}
	// TODO test grid for walls, targets and boxes

}
