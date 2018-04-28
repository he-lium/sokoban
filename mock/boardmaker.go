package mock

import "github.com/he-lium/sokoban"

// Provide mock implementations of BoardMakers that generate hard-coded Boards
type BoardMaker1 struct{}
type BoardMaker2 struct{}
type BoardMaker3 struct{}

var _ sokoban.BoardMaker = (*BoardMaker1)(nil)
var _ sokoban.BoardMaker = (*BoardMaker2)(nil)
var _ sokoban.BoardMaker = (*BoardMaker3)(nil)

func (c BoardMaker1) GenBoard() (*sokoban.Board, error) {
	/* Board structure
	###
	#P#
	###
	*/
	g := sokoban.NewEmptyBoard(0, 3, 3)
	g.AddWall(0, 0)
	g.AddWall(0, 1)
	g.AddWall(0, 2)
	g.AddWall(1, 0)
	g.AddWall(1, 2)
	g.AddWall(2, 0)
	g.AddWall(2, 1)
	g.AddWall(2, 2)
	g.InitPlayer(1, 1)
	return g, nil
}

func (m BoardMaker2) GenBoard() (*sokoban.Board, error) {
	/* Structure
	####
	#P #
	#  #
	####
	*/
	g := sokoban.NewEmptyBoard(1, 4, 4)
	for i := 0; i < 4; i++ {
		g.AddWall(i, 0)
		g.AddWall(i, 3)
	}
	g.AddWall(0, 1)
	g.AddWall(3, 1)
	g.AddWall(0, 2)
	g.AddWall(3, 2)

	g.InitPlayer(1, 1)
	return g, nil
}

func (m BoardMaker3) GenBoard() (*sokoban.Board, error) {
	g := sokoban.NewEmptyBoard(2, 6, 5)
	for i := 0; i < 6; i++ { // top and bottom wall
		g.AddWall(i, 0)
		g.AddWall(i, 4)
	}
	for i := 0; i < 5; i++ { // left and right wall
		g.AddWall(0, i)
		g.AddWall(5, i)
	}
	g.AddWall(3, 2)
	g.AddBox(3, 3)
	g.AddTarget(4, 3)
	g.InitPlayer(4, 2)

	return g, nil
	/* Structure:
	######
	#    #
	#  #P#
	#  BT#
	######
	*/
}
