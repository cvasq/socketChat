package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Message struct {
	Username string `json:"username"`
	Time     string `json:"time"`
	Message  string `json:"message"`
}

type SocketChat struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
}

func newSocketChat() *SocketChat {
	return &SocketChat{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
	}
}

var wsUpgrader = websocket.Upgrader{}

func (h *SocketChat) websocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the initial request to a Websocket
	websocket, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer websocket.Close()

	// Add to client map
	h.clients[websocket] = true

	for {
		var msg Message
		// Accept JSON mapped to Message struct
		err := websocket.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(h.clients, websocket)
			break
		}
		// Send the newly received message to the broadcast channel
		fmt.Println(msg)
		h.broadcast <- msg
	}
}

func (h *SocketChat) handleMessages() {
	for {
		// Blocks until we receive a message on broadcast channel
		msg := <-h.broadcast
		// Send message to all clients
		for client := range h.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Client Write Error: %v", err)
				client.Close()
				delete(h.clients, client)
			}
		}
	}
}
