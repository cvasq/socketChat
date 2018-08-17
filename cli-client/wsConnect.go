package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jroimartin/gocui"
)

var (
	connection net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
)

type Message struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Time     string `json:"time"`
	Data     string `json:"data"`
}

// Disconnect from chat and close
func Disconnect(g *gocui.Gui, v *gocui.View) error {
	connection.Close()
	return gocui.ErrQuit
}

var send = make(chan string)

// Send message
func Send(g *gocui.Gui, v *gocui.View) error {
	send <- v.Buffer()
	g.Update(func(g *gocui.Gui) error {
		v.Clear()
		v.SetCursor(0, 0)
		v.SetOrigin(0, 0)
		return nil
	})
	return nil
}

// Connect to the server, create new reader, writer set client name
func Connect(g *gocui.Gui) error {

	u := url.URL{Scheme: "ws", Host: "localhost", Path: "/ws"}
	s, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	listen := make(chan *Message)

	go func() {
		for {
			receivedData := &Message{}
			err := s.ReadJSON(receivedData)
			if err != nil {
				log.Println("read:", err)
				return
			}
			// Send the newly received message to the broadcast channel
			listen <- receivedData
		}
	}()

	currentTime := func() string {
		const layout = "Jan 2 - 3:04pm"
		now := time.Now()
		return fmt.Sprintf(now.Format(layout))
	}

	subscribeLiveTransactions := func() {
		opMessage := `{"username":"BTC Bot","time":"` + currentTime() + `","data":"Hello, the current price of BTC is $200.55"}`
		err := s.WriteMessage(websocket.TextMessage, []byte(opMessage))
		if err != nil {
			log.Println("Error subscribing to upstream socket:", err)
			s, _, _ = websocket.DefaultDialer.Dial(u.String(), nil)
			return
		}
	}
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for _ = range ticker.C {
			subscribeLiveTransactions()
		}
	}()

	// Some UI changes
	g.SetViewOnTop("intro")
	time.Sleep(time.Second * 3)

	g.SetViewOnTop("messages")
	g.SetViewOnTop("users")
	g.SetViewOnTop("input")
	g.SetCurrentView("input")
	// Wait for server messages in new goroutine
	messagesView, _ := g.View("messages")
	usersView, _ := g.View("users")
	done := make(chan interface{})

	go func() {
	loop:
		for {
			select {
			case msg := <-listen:
				switch {
				case msg.Type == "user-count":
					g.Update(func(g *gocui.Gui) error {
						usersView.Title = fmt.Sprintf(" %v users: ", msg.Data)
						return nil
					})
				case msg.Type == "new-user":
					g.Update(func(g *gocui.Gui) error {
						//usersView.Title = fmt.Sprintf(" %d users: ", clientsCount)
						usersView.Clear()
						fmt.Fprintln(usersView, msg.Username)
						return nil
					})

				default:

					g.Update(func(g *gocui.Gui) error {
						fmt.Fprintln(messagesView, msg)
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
