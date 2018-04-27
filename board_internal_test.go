package sokoban

import "testing"

func TestNewEmptyBoard(t *testing.T) {
	b := NewEmptyBoard(1, 10, 15)
	if b == nil {
		t.Errorf("NewEmptyBoard: returned nil")
		return
	}
	if len(b.Grid) != 10 {
		t.Errorf("NewEmptyBoard: grid width, got: %d, want: %d.", len(b.Grid), 10)
	}
	if len(b.Grid[0]) != 15 {
		t.Errorf("NewEmptyBoard: grid height, got: %d, want: %d.", len(b.Grid[0]), 15)
	}
	b.InitPlayer(5, 7)

	assertPlayer(b, t, 5, 7)
	assertMove(b, t, Up, true)
	assertPlayer(b, t, 5, 6)
	assertMove(b, t, Down, true)
	assertPlayer(b, t, 5, 7)
	assertMove(b, t, Left, true)
	assertPlayer(b, t, 4, 7)
	assertMove(b, t, Right, true)
	assertPlayer(b, t, 5, 7)

	if b.UndoMove() == false {
		t.Error("Undo should return true")
	}
	assertPlayer(b, t, 4, 7)
	if len(b.history) != 3 {
		t.Errorf("history: %d turns, expected 3", len(b.history))
	}

	b.Reset()
	if len(b.history) != 0 {
		t.Errorf("history: %d turns, expected 3", len(b.history))
	}
	assertPlayer(b, t, 5, 7)

}

func assertMove(b *Board, t *testing.T, d Direction, expect bool) {
	if b.MakeMove(d) != expect {
		t.Errorf("MakeMove(%s) should return %t", DirectionToStr(d), expect)
	}
}

func assertPlayer(b *Board, t *testing.T, x, y int) {
	if b.Player.X != x || b.Player.Y != y {
		t.Errorf("player at: (%d, %d), expected: (%d, %d)",
			b.Player.X, b.Player.Y, x, y)
	}
}
