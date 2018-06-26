package websocket

import (
	"log"

	"github.com/he-lium/sokoban"
	"github.com/he-lium/sokoban/parse"
)

// Controller implements sokoban.Controller by communicating to user(s) via
// Go channels
type Controller struct {
	receiver  chan receiveInfo // channel where users send their inputs
	sender    []*client        // slice of channels to send the results
	nPlaying  int              // number of players who haven't left the game
	connected []bool           // bit table of players connected to server
	won       []bool           // bit table of players who have won
}

type receiveInfo struct {
	player int
	json   map[string]interface{}
}

var _ sokoban.Controller = (*Controller)(nil)

// Init broadcasts the initial game board to each user
func (c *Controller) Init(b *sokoban.Board) {
	for i := range c.sender {
		j, err := parse.InitBoardJSON(c.nPlaying, i, b)
		if err == nil {
			c.sendTo(i, j)
		} else {
			// TODO Log error
		}
	}
}

// RecvInput receives an action from the user
// precondition: "action" is a valid json string
func (c *Controller) RecvInput() (int, sokoban.Action) {
	req := <-c.receiver
	var a sokoban.Action
	action := req.json["action"].(string)

	switch action {
	case "undo":
		a.Type = sokoban.Undo
	case "reset":
		a.Type = sokoban.Reset
	case "move":
		a.Type = sokoban.Move
		var direction = ""
		dirField := req.json["direction"]
		if dirField != nil {
			d, ok := dirField.(string)
			if ok {
				direction = d
			}
		}

		switch direction {
		case "up":
			a.Direction = sokoban.Up
		case "down":
			a.Direction = sokoban.Down
		case "left":
			a.Direction = sokoban.Left
		case "right":
			a.Direction = sokoban.Right
		default: // invalid direction
			a.Direction = -1
		}
	case "disconnect":
		// controller action: player has disconnected
		if c.connected[req.player] {
			c.connected[req.player] = false
			close(c.sender[req.player].sendMsg)
			if !c.won[req.player] {
				c.nPlaying--
			}
		}
	default:
		// invalid action
	}

	log.Printf("game %p controller: player %d %s %s",
		c, req.player, action, sokoban.DirectionToStr(a.Direction))
	return req.player, a
}

// SendResult sends the result of an action to user making the action
// and, if successful, broadcasts to all players
func (c *Controller) SendResult(player int, success bool, a sokoban.Action) {
	// send result to the origin player
	c.sendTo(player, parse.ActionResult(player, success))
	if success {
		for i := range c.sender {
			if i != player {
				c.sendTo(i, parse.OpponentAction(player, a))
			}
		}
	}
}

// OutputBoard broadcasts game winners
func (c *Controller) OutputBoard(player int, b *sokoban.Board) {
	// TODO
	if b.Won() {
		c.won[player] = true
		for i := range c.sender {
			c.sendTo(i, parse.WinResultJSON(player))
		}
		c.nPlaying--
	}
}

func (c *Controller) sendTo(p int, msg []byte) {
	// attempt to send to sender goroutine, disconnecting if failed
	if c.connected[p] {
		select {
		case c.sender[p].sendMsg <- msg:
		default:
			log.Printf("controller %p: player %d unresponsive", c, p)
			c.connected[p] = false
			close(c.sender[p].sendMsg)
			if !c.won[p] {
				c.nPlaying--
			}
		}
	}
}

// Closing returns whether the game should stop
func (c *Controller) Closing() bool {
	if c.nPlaying <= 0 {
		// TODO
	}
	return c.nPlaying <= 0
	// After game has ended, caller is responsible for
	// handling channels and Websocket connections
}
