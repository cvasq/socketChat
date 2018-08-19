package main

import (
	"github.com/jroimartin/gocui"
)

// Layout defines the UI parameters
func Layout(g *gocui.Gui) error {

	g.Cursor = true
	g.FgColor = gocui.ColorGreen
	maxX, maxY := g.Size()

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

	if messages, err := g.SetView("messages", 21, 0, maxX-41, maxY-7); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		messages.Title = " Chat Log "
		messages.Autoscroll = true
		messages.Wrap = true
	}

	if input, err := g.SetView("input", 21, maxY-6, maxX-41, maxY-2); err != nil {
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
