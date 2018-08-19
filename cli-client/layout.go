package main

import (
	"fmt"
	"io/ioutil"

	"github.com/jroimartin/gocui"
)

func introBanner() string {
	b, err := ioutil.ReadFile("banner.txt")
	if err != nil {
		panic(err)
	}
	return string(b)
}

// Layout creates chat ui
func Layout(g *gocui.Gui) error {

	g.Cursor = true
	g.FgColor = gocui.ColorGreen
	maxX, maxY := g.Size()

	if intro, err := g.SetView("intro", 0, 0, maxX-40, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		intro.Autoscroll = true
		intro.Wrap = true
		fmt.Fprintf(intro, "\n\n\n\n\n\n\n\t\t\t\t%v", introBanner())
		fmt.Fprintf(intro, "\tGenerating random username...")
	}

	if users, err := g.SetView("users", 0, 0, 20, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		users.Title = " Users Online "
		users.Autoscroll = false
		users.Wrap = true
	}

	if users, err := g.SetView("bots", 0, maxY/2, 20, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		users.Title = " Bots "
		users.Autoscroll = false
		users.Wrap = true
	}

	if messages, err := g.SetView("messages", 21, 0, maxX-40, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		messages.Title = " Chat Log "
		messages.Autoscroll = true
		messages.Wrap = true
	}

	if input, err := g.SetView("input", 21, maxY-6, maxX-40, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		input.Title = " Type message (Press Enter to send): "
		input.Editable = true
		input.Wrap = true
		input.BgColor = gocui.ColorWhite
		input.FgColor = gocui.ColorBlack
	}

	return nil
}
