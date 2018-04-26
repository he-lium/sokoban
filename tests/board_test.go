package tests

import (
	"testing"

	"github.com/he-lium/sokoban"
)

type CellBoard struct{}
type SimpleBoard struct{}
type SimpleBoard2 struct{}

var _ sokoban.BoardMaker = (*CellBoard)(nil)
var _ sokoban.BoardMaker = (*SimpleBoard)(nil)
var _ sokoban.BoardMaker = (*SimpleBoard2)(nil)

func (c CellBoard) GenBoard() *sokoban.Board {
	/* Board structure
	###
	#P#
	###
	*/
	g := sokoban.NewEmptyBoard(0, 3, 3)
	g.AddWall(0, 0)
	g.AddWall(0, 1)
	g.AddWall(0, 2)
	g.AddWall(1, 0)
	g.AddWall(1, 2)
	g.AddWall(2, 0)
	g.AddWall(2, 1)
	g.AddWall(2, 2)
	g.InitPlayer(1, 1)
	return g
}

func TestCellBoard(t *testing.T) {
	c := CellBoard{}
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

func (m SimpleBoard) GenBoard() *sokoban.Board {
	/* Structure
	####
	#P #
	#  #
	####
	*/
	g := sokoban.NewEmptyBoard(1, 4, 4)
	for i := 0; i < 4; i++ {
		g.AddWall(i, 0)
		g.AddWall(i, 3)
	}
	g.AddWall(0, 1)
	g.AddWall(3, 1)
	g.AddWall(0, 2)
	g.AddWall(3, 2)

	g.InitPlayer(1, 1)
	return g
}

func TestMovePlayer(t *testing.T) {
	g := SimpleBoard{}.GenBoard()
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

func (m SimpleBoard2) GenBoard() *sokoban.Board {
	g := sokoban.NewEmptyBoard(2, 6, 5)
	for i := 0; i < 6; i++ { // top and bottom wall
		g.AddWall(i, 0)
		g.AddWall(i, 4)
	}
	for i := 0; i < 5; i++ { // left and right wall
		g.AddWall(0, i)
		g.AddWall(5, i)
	}
	g.AddWall(3, 2)
	g.AddBox(3, 3)
	g.AddTarget(4, 3)
	g.InitPlayer(4, 2)

	return g
	/* Structure:
	######
	#    #
	#  #P#
	#  BT#
	######
	*/
}

func TestWithBox(t *testing.T) {
	g := SimpleBoard2{}.GenBoard()
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
