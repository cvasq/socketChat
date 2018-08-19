package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

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
	var err error
	select {
	case <-done:
		return err
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
func Connect(g *gocui.Gui) error {

	u := url.URL{Scheme: "ws", Host: "localhost", Path: "/ws"}
	s, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error dialing:", err)
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
			case m := <-send:
				err := s.WriteJSON(m)
				if err != nil {
					log.Println("Error sending message:", err)
					return
				}
			default:
			}
		}
	}()

	// Some UI changes
	g.SetViewOnTop("intro")
	messagesView, _ := g.View("messages")
	usersView, _ := g.View("users")

	time.Sleep(3 * time.Second)
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

				case msg.Type == "user-enter":
					g.Update(func(g *gocui.Gui) error {
						fmt.Fprintln(messagesView, msg.Data)
						return nil
					})

				default:

					g.Update(func(g *gocui.Gui) error {
						fmt.Fprintln(messagesView, fmt.Sprintf("%v: %v", msg.Username, msg.Data))
						return nil
					})
				}
			case <-done:
				break loop
			default:
			}
		}
	}()
	return nil
}

// Disconnect from chat and close
func Disconnect(g *gocui.Gui, v *gocui.View) error {
	close(listen)
	close(done)
	return gocui.ErrQuit
}
