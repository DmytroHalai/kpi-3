package painter

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"
	"time"

	"golang.org/x/exp/shiny/screen"
)

var testSize = image.Pt(400, 400)

func TestLoop_Behavior(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)
	l.Receiver = &tr

	var callOrder []string

	// Стартуємо цикл
	l.Start(mockScreen{})

	// Постимо операції, що оновлюють фон
	scene := &Scene{}

	l.Post(OperationFunc(func(tx screen.Texture) {
		callOrder = append(callOrder, "white fill")
		WhiteFill(scene).Do(tx)
	}))

	l.Post(OperationFunc(func(tx screen.Texture) {
		callOrder = append(callOrder, "green fill")
		GreenFill(scene).Do(tx)
	}))

	// Операція, яка не змінює текстуру, але просить її оновити
	l.Post(UpdateOp)

	// Операції з вкладеними Post
	l.Post(OperationFunc(func(screen.Texture) {
		callOrder = append(callOrder, "op 1")
		l.Post(OperationFunc(func(screen.Texture) {
			callOrder = append(callOrder, "op 2")
		}))
	}))
	l.Post(OperationFunc(func(screen.Texture) {
		callOrder = append(callOrder, "op 3")
	}))
	time.Sleep(100 * time.Millisecond) // щоб Loop встиг обробиит op 2

	l.StopAndWait()

	// Перевірка, що текстура була оновлена
	if tr.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	mt, ok := tr.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture type:", tr.lastTexture)
	}
	if mt.Colors[0] != color.White {
		t.Error("First fill color is not white:", mt.Colors)
	}
	if len(mt.Colors) < 2 {
		t.Error("Expected at least 2 colors filled, got:", mt.Colors)
	}

	expectedOrder := []string{"white fill", "green fill", "op 1", "op 3", "op 2"}
	if !reflect.DeepEqual(callOrder, expectedOrder) {
		t.Errorf("Unexpected call order.\nExpected: %v\nGot: %v", expectedOrder, callOrder)
	}
}

type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("not implemented")
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return new(mockTexture), nil
}

func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("not implemented")
}

type mockTexture struct {
	Colors []color.Color
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return testSize }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.Size()}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}

func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Colors = append(m.Colors, src)
}
