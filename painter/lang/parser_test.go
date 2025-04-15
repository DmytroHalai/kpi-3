package lang

import (
	"strings"
	"testing"

	"github.com/DmytroHalai/kpi-3/painter"
)

func TestParser_Parse_White(t *testing.T) {
	input := "white\n"
	parser := &Parser{}
	scene := &painter.Scene{}
	expectedOperations := []painter.Operation{painter.WhiteFill(scene)}

	operations, err := parser.Parse(strings.NewReader(input), scene)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(operations) != len(expectedOperations) {
		t.Fatalf("expected %d operations, got %d", len(expectedOperations), len(operations))
	}
}

func TestParser_Parse_Green(t *testing.T) {
	input := "green\n"
	parser := &Parser{}
	scene := &painter.Scene{}
	expectedOperations := []painter.Operation{painter.GreenFill(scene)}

	operations, err := parser.Parse(strings.NewReader(input), scene)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(operations) != len(expectedOperations) {
		t.Fatalf("expected %d operations, got %d", len(expectedOperations), len(operations))
	}
}

func TestParser_Parse_Update(t *testing.T) {
	input := "update\n"
	parser := &Parser{}
	scene := &painter.Scene{}
	expectedOperations := []painter.Operation{painter.UpdateOp}

	operations, err := parser.Parse(strings.NewReader(input), scene)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(operations) != len(expectedOperations) {
		t.Fatalf("expected %d operations, got %d", len(expectedOperations), len(operations))
	}
}

func TestParser_Parse_BgRect(t *testing.T) {
	input := "bgrect 0 0 10 10\n"
	parser := &Parser{}
	scene := &painter.Scene{}
	expectedOperations := []painter.Operation{painter.BgRectOp(scene, 0, 0, 10, 10)}

	operations, err := parser.Parse(strings.NewReader(input), scene)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(operations) != len(expectedOperations) {
		t.Fatalf("expected %d operations, got %d", len(expectedOperations), len(operations))
	}
}

func TestParser_Parse_Figure(t *testing.T) {
	input := "figure 5 5\n"
	parser := &Parser{}
	scene := &painter.Scene{}
	expectedOperations := []painter.Operation{painter.ShapeOp(scene, 5, 5)}

	operations, err := parser.Parse(strings.NewReader(input), scene)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(operations) != len(expectedOperations) {
		t.Fatalf("expected %d operations, got %d", len(expectedOperations), len(operations))
	}
}

func TestParser_Parse_Move(t *testing.T) {
	input := "move 10 15\n"
	parser := &Parser{}
	scene := &painter.Scene{}
	expectedOperations := []painter.Operation{painter.MoveOp(scene, 10, 15)}

	operations, err := parser.Parse(strings.NewReader(input), scene)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(operations) != len(expectedOperations) {
		t.Fatalf("expected %d operations, got %d", len(expectedOperations), len(operations))
	}
}

func TestParser_Parse_Reset(t *testing.T) {
	input := "reset\n"
	parser := &Parser{}
	scene := &painter.Scene{}
	expectedOperations := []painter.Operation{painter.ResetOp(scene)}

	operations, err := parser.Parse(strings.NewReader(input), scene)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(operations) != len(expectedOperations) {
		t.Fatalf("expected %d operations, got %d", len(expectedOperations), len(operations))
	}
}

func TestParser_Parse_MultipleCommands(t *testing.T) {
	input := `
		white
		move 1 2
		figure 3 4
		update
	`

	parser := &Parser{}
	scene := &painter.Scene{}

	operations, err := parser.Parse(strings.NewReader(input), scene)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(operations) != 4 {
		t.Fatalf("expected 4 operations, got %d", len(operations))
	}
}


func TestParser_Parse_InvalidCommand(t *testing.T) {
	input := "invalid 1 2\n"
	parser := &Parser{}
	scene := &painter.Scene{}

	_, err := parser.Parse(strings.NewReader(input), scene)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestParser_Parse_InvalidArguments(t *testing.T) {
	input := "bgrect 0 0\n"
	parser := &Parser{}
	scene := &painter.Scene{}

	_, err := parser.Parse(strings.NewReader(input), scene)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	input = "figure x y\n"
	_, err = parser.Parse(strings.NewReader(input), scene)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}
