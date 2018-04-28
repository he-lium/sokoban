package main

import (
	"fmt"
	"os"

	"github.com/he-lium/sokoban"
	"github.com/he-lium/sokoban/mock"
	"github.com/he-lium/sokoban/terminal"
)

// Play a single player sokoban game where the terminal displays the board and
// the user enters actions through keyboard characters

func main() {
	gen := mock.BoardMaker3{}
	controller := &terminal.Controller{
		R: os.Stdin,
		W: os.Stdout,
	}
	game, err := sokoban.InitGame(1, gen, controller)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while starting game: %s\n", err.Error())
		os.Exit(2)
	}
	game.Play()
}
