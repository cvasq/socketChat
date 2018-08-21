package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/jroimartin/gocui"
)

type Message struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Time     string `json:"time"`
	Data     string `json:"data"`
}

var send = make(chan Message)
var listen = make(chan *Message)
var done = make(chan interface{})

// Send message on pressing enter key
func Send(g *gocui.Gui, v *gocui.View) error {
	message := Message{
		Type: "user-message",
		Data: strings.TrimSuffix(v.Buffer(), "\n"),
	}

	select {
	case <-done:
		return nil
	default:
		send <- message
	}
	g.Update(func(g *gocui.Gui) error {
		v.Clear()
		v.SetCursor(0, 0)
		v.SetOrigin(0, 0)
		return nil
	})
	return nil
}

// Connect to the Websocket server and start sending/receiving messages
func Connect(g *gocui.Gui, socketchatServerURL *string) error {

	u := url.URL{Scheme: "ws", Host: *socketchatServerURL, Path: "/ws"}
	s, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Unable to establish connection:", err)
	}

	// Receive messages from server
	go func() {
		for {
			receivedData := &Message{}
			err := s.ReadJSON(receivedData)
			if err != nil {
				log.Println("Error receiving message:", err)
				return
			}
			listen <- receivedData
		}
	}()

	// Sends user message to Websocket Connection
	go func() {
		for {
			select {
			case <-done:
				break
			case m := <-send:
				if m.Data != "" {
					err := s.WriteJSON(m)
					if err != nil {
						log.Println("Error sending message:", err)
						return
					}
				}

			}
		}
	}()

	messagesView, _ := g.View("messages")
	usersView, _ := g.View("users")
	botView, _ := g.View("bots")

	g.SetViewOnTop("messages")
	g.SetViewOnTop("users")
	g.SetViewOnTop("input")
	g.SetViewOnTop("bots")
	g.SetCurrentView("input")

	// Wait for server messages in new goroutine
	go func() {
	loop:
		for {
			select {
			case msg := <-listen:
				switch {
				//Client list managed server side
				case msg.Type == "client-list":
					g.Update(func(g *gocui.Gui) error {
						usersView.Title = fmt.Sprintf(" %v User(s) Online ",
							len(strings.Fields(msg.Data)))
						usersView.Clear()
						fmt.Fprintln(usersView, "\n"+msg.Data)
						return nil
					})
				case msg.Type == "bot-list":
					g.Update(func(g *gocui.Gui) error {
						botView.Clear()
						fmt.Fprintln(botView, "\n"+msg.Data)
						return nil
					})
				case msg.Type == "user-enter":
					g.Update(func(g *gocui.Gui) error {
						fmt.Fprintln(messagesView, fmt.Sprintf("\u001b[36;1m[%v] %v\033[0m", msg.Time, msg.Data))
						return nil
					})
				default:
					g.Update(func(g *gocui.Gui) error {
						fmt.Fprintln(messagesView, fmt.Sprintf("[%v] %v: %v", msg.Time, msg.Username, msg.Data))
						return nil
					})
				}
			case <-done:
				break loop
			}
		}
	}()
	return nil
}

// Disconnect from chat and close
func Disconnect(g *gocui.Gui, v *gocui.View) error {
	close(done)
	return gocui.ErrQuit
}
