package main

import (
	"net/http"

	"github.com/DmytroHalai/kpi-3/painter"
	"github.com/DmytroHalai/kpi-3/painter/lang"
	"github.com/DmytroHalai/kpi-3/ui"
)

func main() {
	var (
		pv ui.Visualizer // Візуалізатор створює вікно та малює у ньому.

		// Потрібні для частини 2.
		opLoop painter.Loop // Цикл обробки команд.
		parser lang.Parser  // Парсер команд.
		scene  painter.Scene
		//scene painter.Scene{}

	)

	//pv.Debug = true
	pv.Title = "Simple painter"

	pv.OnScreenReady = opLoop.Start
	opLoop.Receiver = &pv

	go func() {
		http.Handle("/", lang.HttpHandler(&opLoop, &parser, &scene))
		_ = http.ListenAndServe("localhost:17000", nil)
	}()

	pv.Main()
	opLoop.StopAndWait()
}
