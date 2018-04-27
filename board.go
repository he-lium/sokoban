package sokoban

// BoardService defines interface for game controller
type BoardService interface {
	ProcessMove(Direction)
}

// Board is a representation of the Sokoban game board, storing
// positions of objects and move history
type Board struct {
	ID     int
	Width  int
	Height int
	Grid   [][]BoardItem

	Player  Point
	boxes   []Point
	targets []Point

	history []move
	score   int
}

// BoardItem shows what is in the current grid spot.
type BoardItem struct {
	ItemType    BoardItemType
	ContainsBox bool
	boxID       int
	targetID    int
}

// BoardItemType is an enum for the type of grid item.
type BoardItemType int

// Space represents nothing at grid spot.
// Box indicates a box at grid spot.
// Target indicates a target (where the box should be pushed into) at grid spot.
// Wall indicates a wall at grid spot.
const (
	Space  BoardItemType = iota
	Target BoardItemType = iota
	Wall   BoardItemType = iota
)

// Point represents the grid co-ords of the board
type Point struct{ X, Y int }

type move struct {
	from    Point
	to      Point
	boxFrom *Point
	boxTo   *Point
}

// Direction represents the direction in which the player attempts to move
type Direction int

// Direction enums
const (
	Up    Direction = iota
	Right Direction = iota
	Down  Direction = iota
	Left  Direction = iota
)

// DirectionToStr returns the string associated with the given Direction.
func DirectionToStr(d Direction) string {
	switch d {
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Left:
		return "Left"
	case Right:
		return "Right"
	default:
		return "?"
	}
}
