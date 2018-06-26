package websocket

import (
	"log"
	"time"

	"github.com/he-lium/sokoban"
)

const maxRoomPlayers = 2

// Hub matches websocket clients to start the game
type Hub struct {
	register   chan *client
	deregister chan *client
	gen        sokoban.BoardMaker // board generator
	waiting    map[*client]bool   // clients waiting to play
}

// NewHub initialises a waiting hub with given BoardMaker for making new games
func NewHub(gen sokoban.BoardMaker) *Hub {
	return &Hub{
		register:   make(chan *client),
		deregister: make(chan *client),
		waiting:    make(map[*client]bool),
		gen:        gen,
	}
}

// Run starts the hub, receiving connections and spinning off games
func (h *Hub) Run() {
	for {
		select {
		case newClient := <-h.register:
			h.waiting[newClient] = true
			log.Printf("client joined hub. now %d players", len(h.waiting))
			if len(h.waiting) >= maxRoomPlayers {
				// enough people have joined; assign Controller and start game
				h.startNewGame()
			}
		case delClient := <-h.deregister: // quit before playing
			if _, ok := h.waiting[delClient]; ok {
				delete(h.waiting, delClient)
				close(delClient.sendMsg)
				log.Printf("client left hub. now %d players", len(h.waiting))
			}
		}
	}
}

func (h *Hub) startNewGame() {
	// enough people have joined; assign Controller and start game
	numPlayers := len(h.waiting)

	ctrl := &Controller{
		receiver:  make(chan receiveInfo, numPlayers+1),
		sender:    make([]*client, numPlayers),
		nPlaying:  numPlayers,
		connected: make([]bool, numPlayers),
		won:       make([]bool, numPlayers),
	}

	i := 0
	// connect client to controller and remove client from hub
	for c := range h.waiting {
		ctrl.sender[i] = c
		ctrl.connected[i] = true

		c.playLock.Lock()
		c.controller = ctrl
		c.playerID = i
		c.playLock.Unlock()
		delete(h.waiting, c)
		i++
	}
	log.Printf("startNewGame: starting new game with %d players\n", numPlayers)
	// handle playing in separate goroutine
	go runGame(ctrl, h.gen)
}

func runGame(c *Controller, gen sokoban.BoardMaker) {
	defer onFinishGame(c)

	game, err := sokoban.InitGame(c.nPlaying, gen, c)
	if err != nil {
		log.Printf("Hub: ERROR when creating game: %s\n", err.Error())
		return
	}
	game.Play()
}

func onFinishGame(c *Controller) {
	log.Printf("Finished game %p. disconnecting players...\n", c)
	// disconnect clients still playing
	for i := range c.sender {
		if c.connected[i] {
			c.connected[i] = false
			close(c.sender[i].sendMsg)
		}
	}
	go receiveDisconnects(c)
	// Wait for clients to send disconnect messages before closing channel
	time.Sleep(pongWait)
	close(c.receiver)
	log.Printf("Disconnected game %p.\n", c)
}

// receive disconnect calls
func receiveDisconnects(c *Controller) {
	for info := range c.receiver {
		if c.connected[info.player] {
			c.connected[info.player] = false
			close(c.sender[info.player].sendMsg)
		}
	}
}
