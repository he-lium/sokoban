package sokoban

// Interface and definitions for the Game object
// Handles and reroutes actions on Boards from multiple players

// Game acts as a controller for player's game board(s)
type Game struct {
	boards  []*Board
	control Controller
}

// Action represents a player's attempt on making a move
type Action struct {
	Type      ActionType
	Direction Direction
}

// ActionType enum for how the player moves
type ActionType int

// Move: make move in Turn.Direction
// Reset: set the board back to starting state
// Undo: delete last move
const (
	Move  ActionType = 1
	Reset ActionType = 2
	Undo  ActionType = 3
)

// Controller is interface for different types of input e.g. console, web
type Controller interface {
	Init(*Board)
	RecvInput() (int, Action)
	SendResult(player int, success bool, a Action)
	OutputBoard(player int, b *Board)
	Closing() bool
}

// InitGame creates a Game instance with given number of players, controller and board generator
func InitGame(nPlayers int, gen BoardMaker, c Controller) *Game {
	g := &Game{
		boards:  make([]*Board, nPlayers),
		control: c,
	}
	g.boards[0] = gen.GenBoard()
	for i := 1; i < nPlayers; i++ {
		g.boards[i] = g.boards[0].Clone()
	}
	return g
}

// Play plays the Game on a loop, invoking Controller interface functions
// until Closing() returns true
func (g *Game) Play() {
	// Broadcast starting board
	g.control.Init(g.boards[0])

	for !g.control.Closing() {
		p, action := g.control.RecvInput()
		var success bool

		switch action.Type {
		case Move:
			success = g.boards[p].MakeMove(action.Direction)
		case Undo:
			success = g.boards[p].UndoMove()
		case Reset:
			g.boards[p].Reset()
			success = true
		default:
			success = false
		}

		g.control.SendResult(p, success, action)
		g.control.OutputBoard(p, g.boards[p])
	}
}
