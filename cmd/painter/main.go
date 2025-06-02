package main

import (
	"log"
	"net/http"

	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/roman-mazur/architecture-lab-3/painter/lang"
	"github.com/roman-mazur/architecture-lab-3/ui"
)

func main() {
	var (
		pv     ui.Visualizer // Візуалізатор створює вікно та малює у ньому.
		opLoop painter.Loop  // Цикл обробки команд.
		parser lang.Parser   // Парсер команд.
	)

	// pv.Debug = true
	pv.Title = "Simple painter"

	pv.OnScreenReady = opLoop.Start
	opLoop.Receiver = &pv

	go func() {
		http.Handle("/", lang.HttpHandler(&opLoop, &parser))
		log.Println("Starting server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("could not start server: %s", err)
		}
	}()

	pv.Main()
	opLoop.StopAndWait()
}
