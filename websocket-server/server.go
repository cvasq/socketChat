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
	Connection websocket.Conn
	Status	bool
}

type SocketChat struct {
	//clients   map[*websocket.Conn]bool
	clients   []*Client
	broadcast chan Message
}

func createClient(w websocket.Conn) *Client {
	client := &Client{
		Username: "user",
		Connection: w,
		Status: true,

	}
}

func createSocketChat() *SocketChat {
	sockerChat := &SocketChat{
		clients:   []*Client,
		broadcast: make(chan Message),
	}
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

	message := Message{}
	message.Type = "new-user"
	message.Username = randomName
	const layout = "Jan 2 - 3:04pm"
	now := time.Now()
	message.Time = fmt.Sprintf(now.Format(layout))
	message.Data = fmt.Sprintf("%v", "registered")

	h.broadcast <- message

	newClient := new(Client)
	newClient.Username = randomName
	newClient.Connection = make(map[*websocket.Conn]bool)
	newClient.Connection[ws] = true
	log.Println("Adding new client! ", newClient.Username)
	h.clients = append(h.clients, *newClient)

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
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		var msg Message
		// Accept JSON mapped to Message struct
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			for _, v := range h.clients {
				log.Println("Removing: ", v.Username)
				delete(v.Connection, ws)
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
				for k := range client.Connection {
					err := k.WriteJSON(msg)
					if err != nil {
						log.Printf("Client Write Error: %v", err)
						k.Close()
						delete(client.Connection, k)
					}
					k.WriteJSON(msg)
				}

			}
		}
	}
}
