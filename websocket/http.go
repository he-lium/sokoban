package websocket

import (
	"log"
	"net/http"

	"github.com/he-lium/sokoban"
)

const listenAddr = ":8080"

// Serve starts the Server for connecting over websocket
func Serve(gen sokoban.BoardMaker) {
	hub := NewHub(gen)

	go hub.Run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWsClient(hub, w, r)
	})
	log.Println("ListenAndServe at ", listenAddr)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	// TODO serve website
	http.Error(w, "Not found", http.StatusNotFound)
}
