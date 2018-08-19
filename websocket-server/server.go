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
	log.Println("Attempting to remove client...")
	h.clients = append(h.clients[:i], h.clients[i+1:]...)
}

var wsUpgrader = websocket.Upgrader{}

func (h *SocketChat) trackActiveClients() {
	go func() {
		for {

			var clientList string
			for _, client := range h.clients {
				clientList += "-" + client.Username + "\n"
			}

			const layout = "Jan 2 - 3:04pm"
			now := time.Now()

			message := Message{
				Type:     "client-list",
				Username: "system",
				Time:     fmt.Sprintf(now.Format(layout)),
				Data:     clientList,
			}

			h.broadcast <- message

			time.Sleep(4 * time.Second)
		}
	}()
}

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

	greeting := Message{}
	greeting.Type = "user-enter"
	greeting.Username = randomName
	const layout = "Jan 2 - 3:04pm"
	now := time.Now()
	greeting.Time = fmt.Sprintf(now.Format(layout))
	greeting.Data = fmt.Sprintf("[+] User %v has entered", randomName)
	h.broadcast <- greeting

	for {
		var msg Message
		// Accept JSON mapped to Message struct
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error reading JSON: %v", err)
			break
		}
		msg.Username = randomName
		// Send the newly received message to the broadcast channel
		//fmt.Println(msg)
		h.broadcast <- msg
	}
}

func (h *SocketChat) handleMessages() {
	for {
		select {
		case msg := <-h.broadcast:
			log.Println("<-h.broadcast received:", msg)

			for i, client := range h.clients {

				err := client.Connection.WriteJSON(msg)
				if err != nil {
					log.Printf("Client Write Error: %v", err)
					client.Connection.Close()
					h.Remove(i)
					break
				}
			}

		}
	}
}
