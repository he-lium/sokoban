package main

import (
	"io/ioutil"
	"log"

	"github.com/he-lium/sokoban/mock"
	"github.com/he-lium/sokoban/parse"
)

// Generates JSON file from board in mock

func main() {
	b, _ := mock.BoardMaker3{}.GenBoard()

	json, err := parse.BoardToJSON(b)
	if err != nil {
		log.Fatal(err)
	}

	// write json to file
	err = ioutil.WriteFile("example.json", json, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
