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

func TestLoop_UpdatesTexture(t *testing.T) {
	var l Loop
	var tr testReceiver
	scene := &Scene{}
	l.Receiver = &tr
	l.Start(mockScreen{})

	l.Post(WhiteFill(scene))
	l.Post(UpdateOp)

	time.Sleep(50 * time.Millisecond)
	l.StopAndWait()

	if tr.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	mt, ok := tr.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture type:", tr.lastTexture)
	}
	if len(mt.Colors) == 0 || mt.Colors[0] != color.White {
		t.Errorf("Expected white fill, got: %+v", mt.Colors)
	}
}

func TestLoop_ProcessesMultipleOps(t *testing.T) {
	var l Loop
	var tr testReceiver
	scene := &Scene{}
	l.Receiver = &tr
	l.Start(mockScreen{})

	l.Post(WhiteFill(scene))
	l.Post(GreenFill(scene))
	l.Post(UpdateOp)

	time.Sleep(50 * time.Millisecond)
	l.StopAndWait()

	mt := tr.lastTexture.(*mockTexture)
	if len(mt.Colors) < 2 {
		t.Errorf("Expected at least 2 fills, got: %+v", mt.Colors)
	}
}

func TestLoop_NestedPost(t *testing.T) {
	var l Loop
	var tr testReceiver
	callOrder := []string{}

	l.Receiver = &tr
	l.Start(mockScreen{})

	l.Post(OperationFunc(func(screen.Texture) {
		callOrder = append(callOrder, "op 1")
		l.Post(OperationFunc(func(screen.Texture) {
			callOrder = append(callOrder, "op 2")
		}))
	}))
	l.Post(OperationFunc(func(screen.Texture) {
		callOrder = append(callOrder, "op 3")
	}))

	time.Sleep(100 * time.Millisecond)
	l.StopAndWait()

	expected := []string{"op 1", "op 3", "op 2"}
	if !reflect.DeepEqual(callOrder, expected) {
		t.Errorf("Unexpected call order.\nExpected: %v\nGot: %v", expected, callOrder)
	}
}

func TestLoop_StopBlocksUntilDone(t *testing.T) {
	var l Loop
	var tr testReceiver
	l.Receiver = &tr
	l.Start(mockScreen{})

	done := make(chan bool)
	go func() {
		l.Post(UpdateOp)
		time.Sleep(20 * time.Millisecond)
		l.StopAndWait()
		done <- true
	}()

	select {
	case <-done:
		// ok
	case <-time.After(200 * time.Millisecond):
		t.Fatal("StopAndWait did not return in time")
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
