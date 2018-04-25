package sokoban

import "time"

// BoardItem shows what is in the current grid spot.
type BoardItem struct {
	itemType    BoardItemType
	containsBox bool
	boxID       int
	targetID    int
}

// BoardItemType is an enum for the type of grid item.
type BoardItemType int

// Space represents nothing at grid spot.
// Player indicates current player is at grid spot.
// Box indicates a box at grid spot.
// Target indicates an endpoint at grid spot.
// TargetBox indicates a box in its endpoint at grid spot.
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

// Board is a representation of the Sokoban game board, storing
// positions of objects and move history
type Board struct {
	ID     int
	width  int
	height int
	grid   [][]BoardItem

	player  Point
	boxes   []Point
	targets []Point

	history   []move
	score     int
	startTime time.Time
}
