// This file details the methods available for the construction of a board

package sokoban

// BoardMaker defines interface for how boards are generated.
// implementations should use the functions defined in this file to generate
// board
type BoardMaker interface {
	GenBoard() (*Board, error)
}

// NewEmptyBoard returns a pointer to a blank board of given dimensions
func NewEmptyBoard(id, w, h int) *Board {
	var b Board
	b.Width = w
	b.Height = h
	b.ID = id

	b.Grid = make([][]BoardItem, w)
	for i := 0; i < w; i++ {
		b.Grid[i] = make([]BoardItem, h)
	}
	b.boxes = make([]Point, 0, 5)
	b.targets = make([]Point, 0, 5)
	b.history = make([]move, 0, 20)
	return &b
}

// AddWall adds a wall to the wall at the given coordinates
func (b *Board) AddWall(x, y int) {
	b.Grid[x][y].ItemType = Wall
}

// AddTarget adds a box to the board at the given coordinates
func (b *Board) AddTarget(x, y int) {
	b.Grid[x][y].ItemType = Target
	b.Grid[x][y].targetID = len(b.targets)
	b.targets = append(b.targets, Point{x, y})
}

// AddBox adds a box to the board at the given coordinates
func (b *Board) AddBox(x, y int) {
	b.Grid[x][y].ContainsBox = true
	b.Grid[x][y].boxID = len(b.boxes)
	b.boxes = append(b.boxes, Point{x, y})
}

// InitPlayer sets the player's initial starting position
func (b *Board) InitPlayer(x, y int) {
	b.Player = Point{x, y}
}

// Clone copies a board that is in its starting position
func (b *Board) Clone() *Board {
	clone := &Board{
		ID:      b.ID,
		Width:   b.Width,
		Height:  b.Height,
		Grid:    gridCopy(b.Grid),
		Player:  b.Player,
		boxes:   make([]Point, len(b.boxes)),
		targets: make([]Point, len(b.targets)),
		history: make([]move, 0, 20),
		score:   b.score,
	}
	copy(clone.boxes, b.boxes)
	copy(clone.targets, b.targets)
	return clone
}

func gridCopy(g [][]BoardItem) [][]BoardItem {
	dup := make([][]BoardItem, len(g))
	for i := range g {
		dup[i] = make([]BoardItem, len(g[i]))
		copy(dup[i], g[i])
	}
	return dup
}
