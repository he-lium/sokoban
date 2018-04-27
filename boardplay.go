package sokoban

// MakeMove attempts to move the player in the given direction.
// returns true if the player was able to move
// returns false if the player can't move e.g. blocked by wall
func (b *Board) MakeMove(dir Direction) bool {

	dx, dy := directionDelta(dir)
	if dx == invalidDir || dy == invalidDir {
		return false
	}

	// next positions of player and box (if applicable)
	var next, next2 Point
	next = Point{b.Player.X + dx, b.Player.Y + dy}
	if !b.validSpace(next) {
		return false
	}

	var valid bool

	// Check whether player move will push box
	nextGrid := &(b.Grid[next.X][next.Y])
	if nextGrid.ContainsBox {
		// determine whether box can be pushed
		next2 = Point{next.X + dx, next.Y + dy}
		if b.validSpace(next2) && !b.Grid[next2.X][next2.Y].ContainsBox {

			b.moveBox(next, next2)

			// move player and update history
			nextMove := move{
				from:    b.Player,
				to:      next,
				boxFrom: &next,
				boxTo:   &next2}
			b.history = append(b.history, nextMove)
			b.Player = next

			valid = true
		} else {
			// Box can't be pushed into wall or another box
			valid = false
		}
	} else {
		// move player and update history
		nextMove := move{b.Player, next, nil, nil}
		b.history = append(b.history, nextMove)
		b.Player = next

		valid = true
	}

	return valid
}

// UndoMove attempts to undo the last move made by the player
// return false if no moves to undo
func (b *Board) UndoMove() bool {
	if len(b.history) == 0 {
		return false
	}

	lastMove := b.history[len(b.history)-1]
	b.history = b.history[:len(b.history)-1]

	b.Player = lastMove.from
	if lastMove.boxFrom != nil && lastMove.boxTo != nil {
		b.moveBox(*lastMove.boxTo, *lastMove.boxFrom)
	}
	return true
}

// Reset the board back to starting state
func (b *Board) Reset() {
	for b.UndoMove() {
	}
}

// GetScore returns the current score of the game.
func (b *Board) GetScore() int {
	return b.score
}

// Won returns whether the player has won the game
func (b *Board) Won() bool {
	return b.score > 0 && b.score == len(b.targets)
}

// moves the box between given Points, updating state and score as needed.
// precondition: points are valid
func (b *Board) moveBox(from, to Point) {
	src := &(b.Grid[from.X][from.Y])
	dest := &(b.Grid[to.X][to.Y])

	dest.ContainsBox = true
	dest.boxID = src.boxID
	b.boxes[dest.boxID] = to

	src.ContainsBox = false
	src.boxID = 0

	if src.ItemType == Target {
		b.score--
	}
	if dest.ItemType == Target {
		b.score++
	}
}

// validSpace returns whether (x,y) is a valid coordinate and not a wall
func (b *Board) validSpace(p Point) bool {
	return p.X >= 0 && p.X < b.Width &&
		p.Y >= 0 && p.Y < b.Height &&
		b.Grid[p.X][p.Y].ItemType != Wall
}

// returns (dx, dy) deltas for Direction enum
func directionDelta(d Direction) (int, int) {
	switch d {
	case Up:
		return 0, -1
	case Down:
		return 0, 1
	case Left:
		return -1, 0
	case Right:
		return 1, 0
	default:
		return invalidDir, invalidDir
	}
}

// Flag for invalid Direction
const invalidDir = 55
