package sokoban_test

import (
	"testing"

	"github.com/he-lium/sokoban"
	"github.com/he-lium/sokoban/mock"
)

func TestCellBoard(t *testing.T) {
	c := mock.BoardMaker1{}
	game := c.GenBoard()
	score := game.GetScore()
	if score != 0 {
		t.Errorf("incorrect score, received %d, expected %d", score, 0)
	}

	if game.MakeMove(sokoban.Up) != false {
		t.Error("MakeMove(Up) should have returned false")
	}
	if game.MakeMove(sokoban.Down) != false {
		t.Error("MakeMove(Down) should have returned false")
	}
	if game.MakeMove(sokoban.Left) != false {
		t.Error("MakeMove(Left) should have returned false")
	}
	if game.MakeMove(sokoban.Right) != false {
		t.Error("MakeMove(Right) should have returned false")
	}
}

func TestMovePlayer(t *testing.T) {
	g := mock.BoardMaker2{}.GenBoard()
	tables := []struct {
		d        sokoban.Direction
		expected bool
	}{
		{sokoban.Up, false},
		{sokoban.Left, false},
		{sokoban.Down, true},
		{sokoban.Down, false},
		{sokoban.Right, true},
		{sokoban.Right, false},
		{sokoban.Up, true},
		{sokoban.Up, false},
		{sokoban.Right, false},
		{sokoban.Left, true},
		{sokoban.Left, false},
	}
	for i, turn := range tables {
		ret := g.MakeMove(turn.d)
		if ret != turn.expected {
			t.Errorf("Turn %d error: Move(%s) should have returned %t",
				i, sokoban.DirectionToStr(turn.d), turn.expected)
		}
	}
}

func TestWithBox(t *testing.T) {
	g := mock.BoardMaker3{}.GenBoard()
	if g.Won() != false {
		t.Error("Won() should return false at start")
	}

	tables := []struct {
		d        sokoban.Direction
		expected bool
	}{
		{sokoban.Left, false},
		{sokoban.Right, false},
		{sokoban.Down, true},
		{sokoban.Left, true},
		{sokoban.Right, true},
		{sokoban.Up, true},
		{sokoban.Up, true},
		{sokoban.Left, true},
		{sokoban.Left, true},
		{sokoban.Down, true},
		{sokoban.Down, false},

		{sokoban.Left, true},
		{sokoban.Down, true},

		{sokoban.Right, true},
		{sokoban.Right, true},
	}
	for i, turn := range tables {
		ret := g.MakeMove(turn.d)
		if ret != turn.expected {
			t.Errorf("Turn %d error: Move(%s) should have returned %t",
				i, sokoban.DirectionToStr(turn.d), turn.expected)
		}
	}

	if g.Won() != true {
		t.Error("Won() should return true at end")
	}
}

func TestUndo(t *testing.T) {
	g := mock.BoardMaker3{}.GenBoard()
	if g.Won() != false {
		t.Error("Won() should return false at start")
	}

	if g.UndoMove() != false {
		t.Error("UndoMove() should return false at start")
	}

	tables := []struct {
		d        sokoban.Direction
		expected bool
	}{
		{sokoban.Up, true},
		{sokoban.Left, true},
		{sokoban.Left, true},
		{sokoban.Down, true},
		{sokoban.Down, true},
		{sokoban.Right, true},
	}
	for i, turn := range tables {
		ret := g.MakeMove(turn.d)
		if ret != turn.expected {
			t.Errorf("Turn %d error: Move(%s) should have returned %t",
				i, sokoban.DirectionToStr(turn.d), turn.expected)
		}
	}
	if g.Won() != true {
		t.Error("Won() should return true after 7 moves")
	}
	// Attempt to undo previous move
	if g.UndoMove() != true {
		t.Error("UndoMove() should return true after Move")
	}
	if g.Won() != false {
		t.Error("Won() should return false after Undo")
	}
	if g.MakeMove(sokoban.Right) != true {
		t.Error("MakeMove(Right) failed after Undo")
	}
	if g.Won() != true {
		t.Error("Won() should return true at end")
	}
}
