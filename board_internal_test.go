package sokoban

import "testing"

func TestNewEmptyBoard(t *testing.T) {
	b := NewEmptyBoard(1, 10, 15)
	if b == nil {
		t.Errorf("NewEmptyBoard: returned nil")
		return
	}
	if len(b.grid) != 10 {
		t.Errorf("NewEmptyBoard: grid width, got: %d, want: %d.", len(b.grid), 10)
	}
	if len(b.grid[0]) != 15 {
		t.Errorf("NewEmptyBoard: grid height, got: %d, want: %d.", len(b.grid[0]), 15)
	}
}
