package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/gorilla/websocket"
)

type Message struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Time     string `json:"time"`
	Data     string `json:"data"`
}

type Client struct {
	Username   string
	Connection *websocket.Conn
	Status     bool
}

type SocketChat struct {
	clients   []*Client
	broadcast chan Message
}

func createClient(w *websocket.Conn, name string) *Client {
	client := &Client{
		Username:   name,
		Connection: w,
		Status:     true,
	}
	return client
}

func createSocketChat() *SocketChat {
	socketChat := &SocketChat{
		clients:   make([]*Client, 0),
		broadcast: make(chan Message),
	}
	return socketChat
}

// Remove disconnected client from chat
func (h *SocketChat) Remove(i int) {
	h.clients = append(h.clients[:i], h.clients[i+1:]...)
}

var wsUpgrader = websocket.Upgrader{}

func (h *SocketChat) websocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the initial request to a Websocket
	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	randomName := randomdata.SillyName()

	newClient := createClient(ws, randomName)

	log.Println("Adding new client! ", newClient.Username)
	h.clients = append(h.clients, newClient)

	go func() {
		for {
			//m := `{"type": "user-count", "username":"system","time":"now","message":"3"}`
			message := Message{}
			message.Type = "user-count"
			message.Username = "system"
			const layout = "Jan 2 - 3:04pm"
			now := time.Now()
			message.Time = fmt.Sprintf(now.Format(layout))
			message.Data = fmt.Sprintf("%v", len(h.clients))
			log.Println("Sedning message")
			h.broadcast <- message
			log.Println("Users", h.clients)
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		var msg Message
		// Accept JSON mapped to Message struct
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			for i, v := range h.clients {
				log.Println("Removing: ", v.Username)
				h.Remove(i)

			}
			break
		}
		// Send the newly received message to the broadcast channel
		//fmt.Println(msg)
		h.broadcast <- msg
	}
}

func (h *SocketChat) handleMessages() {
	for {
		select {
		case msg := <-h.broadcast:
			for _, client := range h.clients {

				log.Println("<-h.broadcast received:", msg)
				err := client.Connection.WriteJSON(msg)
				if err != nil {
					log.Printf("Client Write Error: %v", err)
					client.Connection.Close()
					break
				}
			}

		}
	}
}
