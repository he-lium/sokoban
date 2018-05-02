package parse

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/he-lium/sokoban"
)

// This file contains objects and functions for converting between Board and its
// JSON format

// board prototype from json
type board struct {
	ID      int     `json:"id"`
	Width   int     `json:"width"`
	Height  int     `json:"height"`
	Player  point   `json:"player"`
	Walls   []point `json:"walls"`
	Targets []point `json:"targets"`
	Boxes   []point `json:"boxes"`
}

// point {x,y}
type point [2]int

// JSONBoard creates sokoban.Board objects from a json obj
type JSONBoard struct {
	JSONContent []byte
}

// ensure sokoban.BoardMaker interface is implemented
var _ sokoban.BoardMaker = (*JSONBoard)(nil)

// GenBoard generates initial board from JSON file
func (gen *JSONBoard) GenBoard() (*sokoban.Board, error) {
	var proto board
	err := json.Unmarshal(gen.JSONContent, &proto)
	if err != nil {
		msg := bytes.Buffer{}
		msg.WriteString("unable to parse json - ")
		msg.WriteString(err.Error())
		return nil, errors.New(msg.String())
	}

	b := sokoban.NewEmptyBoard(proto.ID, proto.Width, proto.Height)
	for _, wall := range proto.Walls {
		b.AddWall(wall[0], wall[1])
	}
	for _, target := range proto.Targets {
		b.AddTarget(target[0], target[1])
	}
	for _, box := range proto.Boxes {
		b.AddBox(box[0], box[1])
	}
	b.InitPlayer(proto.Player[0], proto.Player[1])

	return b, nil
}

// BoardToJSON exports an initial sokoban.Board to json format
func BoardToJSON(game *sokoban.Board) ([]byte, error) {
	return json.Marshal(convertFromBoard(game))
}

func convertFromBoard(game *sokoban.Board) board {
	b := board{
		ID:      game.ID,
		Width:   game.Width,
		Height:  game.Height,
		Player:  point{game.Player.X, game.Player.Y},
		Walls:   make([]point, 0),
		Targets: make([]point, 0),
		Boxes:   make([]point, 0),
	}

	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			if game.Grid[x][y].ItemType == sokoban.Wall {
				b.Walls = append(b.Walls, point{x, y})
			}
			if game.Grid[x][y].ItemType == sokoban.Target {
				b.Targets = append(b.Targets, point{x, y})
			}
			if game.Grid[x][y].ContainsBox {
				b.Boxes = append(b.Boxes, point{x, y})
			}
		}
	}
	return b
}
