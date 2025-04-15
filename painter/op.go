package painter

import (
	"github.com/DmytroHalai/kpi-3/ui"
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

type Shape struct {
	X int
	Y int
}

type Rectangle struct {
	X1 int
	Y1 int
	X2 int
	Y2 int
}

type Scene struct {
	BgColor color.Color
	Rect    *Rectangle
	Shapes  []Shape
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

func render(scene *Scene, t screen.Texture) {
	bgColor := scene.BgColor
	if bgColor == nil {
		bgColor = color.RGBA{G: 128, A: 255}
	}
	t.Fill(t.Bounds(), bgColor, screen.Src)
	rect := scene.Rect
	if rect != nil {
		t.Fill(image.Rect(rect.X1, rect.Y1, rect.X2, rect.Y2), color.Black, screen.Src)
	}
	for _, shape := range scene.Shapes {
		ui.DrawTShape(t, shape.X, shape.Y, t.Bounds(), color.RGBA{255, 255, 0, 255})
	}
}

func WhiteFill(scene *Scene) Operation {
	return OperationFunc(func(t screen.Texture) {
		scene.BgColor = color.White
		render(scene, t)
	})
}

func GreenFill(scene *Scene) Operation {
	return OperationFunc(func(t screen.Texture) {
		scene.BgColor = color.RGBA{G: 255, A: 1}
		render(scene, t)
	})
}

func BgRectOp(scene *Scene, x1, y1, x2, y2 int) Operation {
	return OperationFunc(func(t screen.Texture) {
		scene.Rect = &Rectangle{x1, y1, x2, y2}
		render(scene, t)
	})
}

func ShapeOp(scene *Scene, x1, x2 int) Operation {
	return OperationFunc(func(t screen.Texture) {
		scene.Shapes = append(scene.Shapes, Shape{x1, x2})
		render(scene, t)
	})
}

func MoveOp(scene *Scene, x, y int) Operation {
	return OperationFunc(func(t screen.Texture) {
		newShapes := make([]Shape, len(scene.Shapes))
		for i := range scene.Shapes {
			newShapes[i] = Shape{x, y}
		}
		scene.Shapes = newShapes
		render(scene, t)
	})
}

func ResetOp(scene *Scene) Operation {
	return OperationFunc(func(t screen.Texture) {
		scene.BgColor = color.White
		scene.Rect = nil
		scene.Shapes = nil
		render(scene, t)
	})
}
