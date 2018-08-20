package main

import (
	"flag"
	"log"

	"github.com/jroimartin/gocui"
)

func main() {

	// Set custom port by running with --port PORT_NUM
	// Default port is 8080
	socketchatServerURL := flag.String("server-address",
		"localhost:9001",
		"Chat host address")

	flag.Parse()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	g.SetManagerFunc(Layout)
	go Connect(g, socketchatServerURL)
	g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, Send)
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Disconnect)
	g.MainLoop()
}
