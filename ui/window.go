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

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz  size.Event
	pos image.Rectangle
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})

	// Початкове положення фігури "Т" по центру
	pw.pos = image.Rect(300, 250, 500, 450)
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
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
		close(pw.done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.w = w

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

		case t = <-pw.tx:
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
		pw.sz = e
		log.Println("Size event:", e)

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if t == nil {
			if e.Button == mouse.ButtonLeft {
				// Оновлення положення фігури "Т"
				pw.pos = image.Rect(int(e.X)-100, int(e.Y)-100, int(e.X)+100, int(e.Y)+100)
				log.Println("Mouse event - position:", pw.pos)
				pw.w.Send(paint.Event{})
			}
		}

	case paint.Event:
		log.Println("Paint event")
		pw.drawDefaultUI()
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI() {
	log.Println("Drawing default UI")
	pw.w.Fill(pw.sz.Bounds(), color.RGBA{0, 255, 0, 255}, draw.Src) // Фон

	// Малювання фігури "Т"
	figColor := color.RGBA{255, 255, 0, 255}

	// Горизонтальна частина "Т"
	horizontalRect := image.Rect(pw.pos.Min.X-25, pw.pos.Min.Y, pw.pos.Max.X+25, pw.pos.Min.Y+100)
	pw.w.Fill(horizontalRect, figColor, draw.Src)

	// Вертикальна частина "Т"
	verticalRect := image.Rect(pw.pos.Min.X+25, pw.pos.Min.Y+50, pw.pos.Max.X-25, pw.pos.Max.Y+50)
	pw.w.Fill(verticalRect, figColor, draw.Src)
}
