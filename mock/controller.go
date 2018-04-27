package mock

import (
	"testing"

	"github.com/he-lium/sokoban"
)

// Controller is a mock implementation of the sokoban.Controller interface.
// Used for testing sokoban.Game implementation
type Controller struct {
	T *testing.T
	// number of interface function calls
	InitInvoked    int
	RecvInvoked    int
	SendInvoked    int
	OutInvoked     int
	ClosingInvoked int
	// Actions to be passed to game
	Actions []sokoban.Action
	// expected Results correlating to each action
	Results []bool
}

// Ensure MockController1 implements interface
var _ sokoban.Controller = (*Controller)(nil)

func (c *Controller) Init(b *sokoban.Board) {
	c.InitInvoked++
	if c.InitInvoked > 1 {
		c.T.Errorf("Init() called %d times", c.InitInvoked)
	}
	if b.Won() {
		c.T.Error("initial board should not be in winning state")
	}
}

func (c *Controller) RecvInput() (int, sokoban.Action) {
	a := c.Actions[c.RecvInvoked]
	c.RecvInvoked++
	return 0, a
}

func (c *Controller) SendResult(p int, success bool, a sokoban.Action) {
	if c.Results[c.SendInvoked] != success {
		c.T.Errorf("result at turn %d should be %t", c.SendInvoked, success)
	}
	c.SendInvoked++
	if c.RecvInvoked != c.SendInvoked {
		c.T.Errorf("RecvInput() called %d times, SendResult() called %d times",
			c.RecvInvoked, c.SendInvoked)
	}
}

func (c *Controller) OutputBoard(p int, b *sokoban.Board) {
	c.OutInvoked++
	if c.SendInvoked != c.OutInvoked {
		c.T.Errorf("SendResult() called %d times, OutputBoard() called %d times",
			c.SendInvoked, c.OutInvoked)
	}
	// if board is on last move, assert that player has won
	if c.OutInvoked == len(c.Results) && !b.Won() {
		c.T.Errorf("Did not win game on last round")
	}
}

func (c *Controller) Closing() bool {
	c.ClosingInvoked++
	return c.ClosingInvoked-1 >= len(c.Results)
}
