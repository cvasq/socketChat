package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
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
	clients    []*Client
	broadcast  chan Message
	triggerBot chan string
}

func currentTime() string {
	const layout = "3:04pm"
	now := time.Now()
	return fmt.Sprintf(now.Format(layout))
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
		clients:    make([]*Client, 0),
		broadcast:  make(chan Message),
		triggerBot: make(chan string),
	}
	return socketChat
}

var wsUpgrader = websocket.Upgrader{}

func (h *SocketChat) websocketHandler(w http.ResponseWriter, r *http.Request) {

	// Upgrade the initial request to a Websocket
	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Move this to new file - newClient.go
	// Genrate random username
	randomName := randomdata.SillyName()
	newClient := createClient(ws, randomName)
	log.Printf("Adding new client %v from %v \n", newClient.Username, ws.RemoteAddr())
	h.clients = append(h.clients, newClient)
	greeting := Message{
		Type:     "user-enter",
		Username: randomName,
		Time:     currentTime(),
		Data:     fmt.Sprintf("++ User %v has entered", randomName),
	}
	h.broadcast <- greeting

	for {

		// Accept JSON mapped to Message struct
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading JSON: %v", err)
			break
		}
		msg.Username = randomName
		msg.Time = currentTime()

		// Send the newly received message to the broadcast channel
		h.broadcast <- msg
	}
}

func (h *SocketChat) handleMessages() {
	// Matches a "/btc" command
	cmd, _ := regexp.Compile("(^[/](btc$))")

	for {
		select {
		case msg := <-h.broadcast:
			//Debug
			log.Println("<-h.broadcast received:", msg)

			// Check whether we received a command
			if cmd.MatchString(msg.Data) == true {

				h.triggerBot <- "btc-bot"

				// Broadcast message if cmd regex does not match
			} else {

				for clientIndex, client := range h.clients {
					err := client.Connection.WriteJSON(msg)
					if err != nil {
						log.Printf("Client Write Error: %v", err)
						client.Connection.Close()
						h.removeClient(clientIndex, client.Username)
						break
					}
				}

			}

		}
	}
}

// Move to new clients.go file
func (h *SocketChat) trackActiveClients() {
	go func() {
		for {

			var clientList string
			for _, client := range h.clients {
				clientList += client.Username + "\n"
			}

			message := Message{
				Type:     "client-list",
				Username: "system",
				Time:     currentTime(),
				Data:     clientList,
			}

			if clientList != "" {
				h.broadcast <- message

				botListing := Message{
					Type:     "bot-list",
					Username: "system",
					Time:     currentTime(),
					Data:     "btc-bot",
				}
				h.broadcast <- botListing
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

// Move to new clients.go file
// Remove disconnected client from chat
func (h *SocketChat) removeClient(i int, client string) {
	log.Println("Removing client:", client)
	h.clients = append(h.clients[:i], h.clients[i+1:]...)
}
