package sokoban_test

import (
	"testing"

	"github.com/he-lium/sokoban"
	"github.com/he-lium/sokoban/mock"
)

func TestGameSinglePlayer(t *testing.T) {
	c := mock.Controller{T: t}
	c.Actions = []sokoban.Action{
		sokoban.Action{Type: sokoban.Move, Direction: sokoban.Right},
		sokoban.Action{Type: sokoban.Move, Direction: sokoban.Up},
		sokoban.Action{Type: sokoban.Move, Direction: sokoban.Left},
		sokoban.Action{Type: sokoban.Move, Direction: sokoban.Left},
		sokoban.Action{Type: sokoban.Move, Direction: sokoban.Down},
		sokoban.Action{Type: sokoban.Move, Direction: sokoban.Down},
		sokoban.Action{Type: sokoban.Move, Direction: sokoban.Right},
		sokoban.Action{Type: sokoban.Undo, Direction: sokoban.Right},
		sokoban.Action{Type: sokoban.Move, Direction: sokoban.Right},
	}
	c.Results = []bool{false, true, true, true, true, true, true, true, true}
	if len(c.Actions) != len(c.Results) {
		panic("each action should have a result")
	}
	g, err := sokoban.InitGame(1, mock.BoardMaker3{}, &c)
	if err != nil {
		c.T.Fatalf("unable to init game: %s", err)
	}

	g.Play()

	if c.InitInvoked != 1 {
		c.T.Errorf("InitInvoked() was called %d times at end", c.InitInvoked)
	}
	if c.RecvInvoked != len(c.Actions) {
		c.T.Errorf("c.RecvInvoked is %d, expected %d", c.RecvInvoked, len(c.Actions))
	}
	if c.SendInvoked != len(c.Results) {
		c.T.Errorf("c.SendInvoked is %d, expected %d", c.SendInvoked, len(c.Results))
	}
}
