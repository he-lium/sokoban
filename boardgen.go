// This file details the methods available for the construction of a board

package sokoban

// BoardMaker defines interface for how boards are generated.
// implementations should use the functions defined in this file to generate
// board
type BoardMaker interface {
	GenBoard() *Board
}

// NewEmptyBoard returns a pointer to a blank board of given dimensions
func NewEmptyBoard(id, w, h int) *Board {
	var b Board
	b.width = w
	b.height = h
	b.ID = id

	b.grid = make([][]BoardItem, w)
	for i := 0; i < w; i++ {
		b.grid[i] = make([]BoardItem, h)
	}
	b.boxes = make([]Point, 0, 5)
	b.targets = make([]Point, 0, 5)
	b.history = make([]move, 0, 20)
	return &b
}

// AddWall adds a wall to the wall at the given coordinates
func (b *Board) AddWall(x, y int) {
	b.grid[x][y].itemType = Wall
}

// AddTarget adds a box to the board at the given coordinates
func (b *Board) AddTarget(x, y int) {
	b.grid[x][y].itemType = Target
	b.grid[x][y].targetID = len(b.targets)
	b.targets = append(b.targets, Point{x, y})
}

// AddBox adds a box to the board at the given coordinates
func (b *Board) AddBox(x, y int) {
	b.grid[x][y].containsBox = true
	b.grid[x][y].boxID = len(b.boxes)
	b.boxes = append(b.boxes, Point{x, y})
}

// InitPlayer sets the player's initial starting position
func (b *Board) InitPlayer(x, y int) {
	b.player = Point{x, y}
}
