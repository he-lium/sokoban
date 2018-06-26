package websocket

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type status int

const (
	// writeWait  = 10 * time.Second  // time allowed to write msg
	// pongWait   = 60 * time.Second  // time allowed to read next pong msg
	writeWait  = 2 * time.Second   // debugging
	pongWait   = 2 * time.Second   // debugging
	pingPeriod = pongWait * 9 / 10 // interval to send pings
	maxMsgSize = 512
)

var newline = []byte{'\n'}
var space = []byte{' '}

// Client is a representation of a websocket connection with a client
type client struct {
	sendMsg    chan []byte // channel of outbound messages
	conn       *websocket.Conn
	hub        *Hub
	controller *Controller
	playerID   int        // assigned at start of game
	isPlaying  bool       // whether the game has started with a Controller
	playLock   sync.Mutex // mutex for initial setup of Controller
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// handle incoming messages from user
func (c *client) incoming() {
	// cleanup
	defer func() {
		if !c.isPlaying {
			c.hub.deregister <- c
		} else {
			// TODO disconnect from Controller if still playing
			j := map[string]interface{}{
				"action": "disconnect",
			}
			c.controller.receiver <- receiveInfo{c.playerID, j}
		}
		log.Printf("player %d client reader closed\n", c.playerID)
		c.conn.Close()
	}()

	// set up
	c.conn.SetReadLimit(maxMsgSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		log.Printf("pong handler for player %d in %p", c.playerID, c.controller)
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		// determine whether still waiting for game to start
		if !c.isPlaying {
			c.playLock.Lock()
			if c.controller != nil {
				log.Printf("Client: started game %p as player %d\n", c.controller, c.playerID)
				c.isPlaying = true
			}
			c.playLock.Unlock()
		}

		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("player %d client read error: %v", c.playerID, err)
			if websocket.IsUnexpectedCloseError(
				err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ws closing error for player %d: %v", c.playerID, err)
			}
			break
		}

		if !c.isPlaying {
			continue
		}
		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))

		log.Printf("client %p:%d: %s\n", c.controller, c.playerID, msg)

		// Parse JSON message, check for action field
		var data map[string]interface{}
		err = json.Unmarshal(msg, &data)
		if err != nil {
			// TODO reply with error msg to client
			continue
		}
		action, ok1 := data["action"]
		_, ok2 := action.(string)
		if !(ok1 && ok2) {
			// TODO reply with error msg to client
			continue
		}

		// send action to game Controller
		c.controller.receiver <- receiveInfo{c.playerID, data}
	}
}

// Process messages to be sent to the peer. Can send one or more JSON objects in
// a message, separated by newline.
func (c *client) outgoing() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		log.Printf("player %d client writer closed\n", c.playerID)
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.sendMsg:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// channel closed
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				// TODO log
				return
			}
			w.Write(msg)

			// Send queued logs to the message
			n := len(c.sendMsg)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.sendMsg)
			}

			if err := w.Close(); err != nil {
				// TODO log
				return
			}

		case <-ticker.C: // time to send ping to peer
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			log.Printf("sending ping msg for player %d in %p", c.playerID, c.controller)
			if err != nil {
				// TODO log
				return
			}
		}
	}
}

func serveWsClient(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("websocket client connected from %s\n", r.RemoteAddr)

	client := &client{
		sendMsg:  make(chan []byte, 5),
		conn:     conn,
		hub:      hub,
		playerID: -1,
	}
	hub.register <- client

	// start message sender and receivers
	go client.incoming()
	go client.outgoing()
}
