package main

import (
	"log"

	"github.com/he-lium/sokoban/mock"
	"github.com/he-lium/sokoban/websocket"
)

func main() {
	log.Println("Sokoban websocket server")
	websocket.Serve(&mock.BoardMaker3{})
}
