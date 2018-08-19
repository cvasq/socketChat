package main

import (
	"github.com/jroimartin/gocui"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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
