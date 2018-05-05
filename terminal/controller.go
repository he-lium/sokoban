package terminal

import (
	"bufio"
	"fmt"
	"io"

	"github.com/he-lium/sokoban"
)

// Controller implements sokoban.Controller and interacts with the user via
// command line, printing the board and accepting keystrokes as moves
type Controller struct {
	R          io.Reader
	W          io.Writer
	valid      bool
	nWon       int
	won        []bool
	NPlayers   int
	currPlayer int
	reader     *bufio.Reader
}

var _ sokoban.Controller = (*Controller)(nil)

// Init prints the initial state of the board to the user
func (c *Controller) Init(b *sokoban.Board) {
	fmt.Fprintln(c.W, "Welcome to 倉庫番!")
	c.reader = bufio.NewReader(c.R)
	showBoard(c.W, b)

	c.won = make([]bool, c.NPlayers)
}

// RecvInput asks the user for an action
func (c *Controller) RecvInput() (int, sokoban.Action) {
	var a sokoban.Action

	for c.won[c.currPlayer] {
		c.currPlayer = (c.currPlayer + 1) % c.NPlayers
	}
	if c.NPlayers > 1 {
		fmt.Fprintf(c.W, "Player %d: ", c.currPlayer+1)
	}

	switch c.prompt() {
	case 'w':
		a = sokoban.Action{Type: sokoban.Move, Direction: sokoban.Up}
	case 'a':
		a = sokoban.Action{Type: sokoban.Move, Direction: sokoban.Left}
	case 's':
		a = sokoban.Action{Type: sokoban.Move, Direction: sokoban.Down}
	case 'd':
		a = sokoban.Action{Type: sokoban.Move, Direction: sokoban.Right}
	case 'u':
		a.Type = sokoban.Undo
	case 'r':
		a.Type = sokoban.Reset
	}
	return c.currPlayer, a
}

func (c *Controller) prompt() rune {
	fmt.Fprintln(c.W, `Select Actions:
(w) Up   (a) Left   (s) Down   (d) Right   (u) Undo   (r) Restart`)
	r, _, err := c.reader.ReadRune()
	for err != nil || r == '\n' {
		fmt.Fprint(c.W, "> ")
		r, _, err = c.reader.ReadRune()
	}
	return r
}

// SendResult shows whether the user's action was successful
func (c *Controller) SendResult(p int, success bool, a sokoban.Action) {
	if !success {
		fmt.Fprintln(c.W, "Invalid action")
	} else {
		c.currPlayer = (c.currPlayer + 1) % c.NPlayers
	}
	c.valid = success
}

// OutputBoard prints the current board
func (c *Controller) OutputBoard(p int, b *sokoban.Board) {
	if c.valid {
		showBoard(c.W, b)
		if b.Won() && !c.won[p] {
			fmt.Fprintln(c.W, "You win!")
			c.won[p] = true
			c.nWon++
		}
	}
}

// Closing signals whether the game has been won
func (c *Controller) Closing() bool {
	return c.nWon == c.NPlayers
}

func showBoard(w io.Writer, b *sokoban.Board) {
	if len(b.Grid) == 0 || len(b.Grid[0]) == 0 {
		return
	}

	for y := range b.Grid[1] {
		for x := range b.Grid {
			if b.Player.X == x && b.Player.Y == y {
				fmt.Fprint(w, "P")
			} else if b.Grid[x][y].ItemType == sokoban.Wall {
				fmt.Fprint(w, "#")
			} else if b.Grid[x][y].ContainsBox {
				fmt.Fprint(w, "B")
			} else if b.Grid[x][y].ItemType == sokoban.Target {
				fmt.Fprint(w, "T")
			} else {
				fmt.Fprint(w, " ")
			}
		}
		fmt.Fprintln(w)
	}
}
