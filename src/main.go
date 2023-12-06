package main

import (
	"log/slog"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/widget/material"
	"github.com/vlabo/vet_note/src/storage"
	"github.com/vlabo/vet_note/src/view"
)

func main() {
	go func() {
		w := app.NewWindow()

		err := run(w)
		if err != nil {
			slog.Error(err.Error())
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var ops op.Ops

	// Get the proper os place to store data.
	dataDir, err := app.DataDir()
	if err != nil {
		return err
	}

	// Initialize storage
	err = storage.InitStorage(dataDir + "/vet_note")
	if err != nil {
		return err
	}
	defer storage.CloseStorage()

	// Initialize theme
	theme := material.NewTheme()
	// theme.ContrastBg = color.NRGBA{R: 230, G: 230, B: 230, A: 255}
	theme.TextSize = 20

	// Views
	selectView := view.NewSelectView(theme)
	patientView := view.NewPatientView(theme)

	var currentView view.View = selectView

	for {
		// Read and process next event
		e := w.NextEvent()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			// Setup keyboard event for the frame.
			area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
			key.InputOp{
				Tag:  w,
				Keys: key.NameEscape,
			}.Add(gtx.Ops)
			for _, event := range gtx.Events(w) {
				switch event := event.(type) {
				case key.Event:
					processKeyboardEvent(w, event)
				}
			}
			currentView.Layout(gtx)
			area.Pop()
			e.Frame(gtx.Ops)
		}

		// Check for view events
		select {
		case e := <-view.ViewEventChan:
			switch e {
			case view.OpenPatientView:
				currentView = patientView
			case view.OpenSelectView:
				currentView = selectView
				selectView.ReloadPatientList()
			}

		default:
			// No View event
		}
	}
}

func processKeyboardEvent(w *app.Window, event key.Event) {
	slog.Info("key", "event", event)
	switch event.Name {
	case key.NameEscape:
		w.Perform(system.ActionClose)
	}
}
