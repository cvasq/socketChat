package main

import (
	"github.com/jroimartin/gocui"
	"log"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	g.SetManagerFunc(Layout)
	go Connect(g)
	g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, Send)
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Disconnect)
	g.MainLoop()
}
