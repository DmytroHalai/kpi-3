package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
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
	pw.pos.Max.X = 200
	pw.pos.Max.Y = 200
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

	events := make(chan any)
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
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {

	case size.Event: // Оновлення даних про розмір вікна.
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if e.Button == mouse.ButtonLeft && e.Direction == mouse.DirPress {
			pw.pos.Min.X = int(e.X)
			pw.pos.Min.Y = int(e.Y)
			pw.w.Send(paint.Event{}) 
		}

	case paint.Event:
		// Малювання контенту вікна.
		if t == nil {
			pw.drawDefaultUI()
		} else {
			// Використання текстури отриманої через виклик Update.
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI() {
	const borderThickness = 10

	bgColor := color.RGBA{0, 128, 0, 255}
	shapeColor := color.RGBA{255, 255, 0, 255}
	borderColor := color.White

	pw.drawBackground(bgColor)
	centerX, centerY := pw.getCenter()
	pw.drawTShape(centerX, centerY, shapeColor)
	pw.drawBorder(borderThickness, borderColor)
}

func (pw *Visualizer) drawBackground(bg color.RGBA) {
	pw.w.Fill(pw.sz.Bounds(), bg, draw.Src)
}

func (pw *Visualizer) getCenter() (int, int) {
	if pw.pos.Min.X != 0 && pw.pos.Min.Y != 0 {
		return pw.pos.Min.X, pw.pos.Min.Y
	}
	return pw.sz.WidthPx / 2, pw.sz.HeightPx / 2
}

func (pw *Visualizer) drawTShape(centerX, centerY int, shapeColor color.RGBA) {
	maxWidth := pw.sz.WidthPx / 2
	maxHeight := pw.sz.HeightPx / 2

	tWidthTop := int(float64(maxWidth) * 0.7)
	tHeightTop := int(float64(maxHeight) * 0.2)
	tWidthVert := int(float64(maxWidth) * 0.2)
	tHeightVert := int(float64(maxHeight) * 0.7)

	topRect := image.Rect(
		centerX - tWidthTop/2,
		centerY - tHeightVert/2,
		centerX + tWidthTop/2,
		centerY - tHeightVert/2 + tHeightTop,
	)

	vertRect := image.Rect(
		centerX - tWidthVert/2,
		centerY - tHeightVert/2,
		centerX + tWidthVert/2,
		centerY + tHeightVert/2,
	)

	pw.w.Fill(vertRect, shapeColor, draw.Src)
	pw.w.Fill(topRect, shapeColor, draw.Src)
}

func (pw *Visualizer) drawBorder(thickness int, borderColor color.Color) {
	for _, br := range imageutil.Border(pw.sz.Bounds(), thickness) {
		pw.w.Fill(br, borderColor, draw.Src)
	}
}