package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	W       screen.Window
	Tx      chan screen.Texture
	Done    chan struct{}
	BgColor color.Color
	BgRect  image.Rectangle

	Sz  size.Event
	Pos image.Rectangle
}

func (pw *Visualizer) Main() {
	pw.Tx = make(chan screen.Texture)
	pw.Done = make(chan struct{})

	// Initializing background color and rect
	pw.BgColor = color.RGBA{0, 255, 0, 255}
	pw.BgRect = image.Rect(0, 0, 0, 0)

	// Initial position of the "T" shape in the center
	pw.Pos = image.Rect(362, 284, 462, 434)
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.Tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title:  pw.Title,
		Width:  800,
		Height: 800,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.Done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.W = w

	events := make(chan interface{})
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.Tx:
			log.Println("Received texture update")
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e interface{}) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e interface{}, t screen.Texture) {
	switch e := e.(type) {

	case size.Event:
		pw.Sz = e
		log.Println("Size event:", e)

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if t == nil {
			if e.Button == mouse.ButtonLeft {
				// Update the position of the "T" shape
				pw.Pos = image.Rect(int(e.X)-50, int(e.Y)-75, int(e.X)+50, int(e.Y)+75)
				log.Println("Mouse event - position:", pw.Pos)
				pw.W.Send(paint.Event{})
			}
		}

	case paint.Event:
		log.Println("Paint event")
		pw.drawDefaultUI()
		pw.W.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI() {
	log.Println("Drawing default UI")
	pw.W.Fill(pw.Sz.Bounds(), pw.BgColor, draw.Src) // Background

	if pw.BgRect != (image.Rectangle{}) {
		pw.W.Fill(pw.BgRect, color.Black, draw.Src)
	}

	// Drawing the "T" shape
	figColor := color.RGBA{255, 255, 0, 255}

	// Horizontal part of the "T"
	horizontalRect := image.Rect(pw.Pos.Min.X-25, pw.Pos.Min.Y, pw.Pos.Max.X+25, pw.Pos.Min.Y+50)
	pw.W.Fill(horizontalRect, figColor, draw.Src)

	// Vertical part of the "T"
	verticalRect := image.Rect(pw.Pos.Min.X+25, pw.Pos.Min.Y+50, pw.Pos.Max.X-25, pw.Pos.Max.Y+50)
	pw.W.Fill(verticalRect, figColor, draw.Src)
}
