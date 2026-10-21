// Author: HyperparameterSI

// Terminal is a terminal app in native UI: the user's shell in a window,
// with the terminal plugin, which emulates it with Ghostty's libghostty-vt.
// The window takes the shell's title and closes when the shell exits.
//
//	go run ./examples/terminal
package main

import (
	"log"

	"github.com/HyperparameterSI/golet/backend"
	"github.com/HyperparameterSI/golet/backend/plugins/terminal"
	"github.com/HyperparameterSI/golet/backend/ui"
)

func main() {
	golet.App.WhenReady(func() {
		var win *golet.Window
		term, err := terminal.New(terminal.Options{
			OnTitle: func(title string) { win.SetTitle(title) },
			OnExit:  func(int) { win.Close() },
		})
		if err != nil {
			log.Fatal(err)
		}
		win = golet.NewWindow(golet.WindowOptions{
			Title:  "Terminal",
			Width:  760,
			Height: 480,
			Content: ui.View(func(c *ui.Context) {
				terminal.View(c, term).Fill().AutoFocus()
			}),
		})
		win.OnClosed(func() { term.Close() })
	})
	if err := golet.App.Run(); err != nil {
		log.Fatal(err)
	}
}
