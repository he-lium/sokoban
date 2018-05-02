package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/he-lium/sokoban"
	"github.com/he-lium/sokoban/parse"
	"github.com/he-lium/sokoban/terminal"
)

// Play a single player sokoban game where the terminal displays the board and
// the user enters actions through keyboard characters

// Board is loaded from json file given in argument
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <board.json>\n", os.Args[0])
		os.Exit(1)
	}
	json, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %s\n", os.Args[1], err.Error())
		os.Exit(2)
	}
	gen := &parse.JSONBoard{JSONContent: json}
	controller := &terminal.Controller{
		R: os.Stdin,
		W: os.Stdout,
	}
	game, err := sokoban.InitGame(1, gen, controller)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while starting game: %s\n", err.Error())
		os.Exit(3)
	}
	game.Play()
}
