package lang

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/DmytroHalai/kpi-3/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)


	for scanner.Scan() {
		commandLine := strings.TrimSpace(scanner.Text())
		if commandLine == "" {
			continue 
		}

		parts := strings.Fields(commandLine)
		if len(parts) == 0 {
			continue 
		}

		instruction := parts[0]
		args := parts[1:]

		switch instruction {
		case "white":
			if len(args) != 0 {
				return nil, fmt.Errorf("white command takes no arguments, got %d", len(args))
			}
			res = append(res, painter.OperationFunc(painter.WhiteFill))

		case "green":
			if len(args) != 0 {
				return nil, fmt.Errorf("green command takes no arguments, got %d", len(args))
			}
			res = append(res, painter.OperationFunc(painter.GreenFill))

		case "update":
			if len(args) != 0 {
				return nil, fmt.Errorf("update command takes no arguments, got %d", len(args))
			}
			res = append(res, painter.UpdateOp)

		case "bgrect":
			if len(args) != 4 {
				return nil, fmt.Errorf("bgrect command requires 4 arguments, got %d", len(args))
			}
			x1, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return nil, fmt.Errorf("bgrect x1 invalid: %v", err)
			}
			y1, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return nil, fmt.Errorf("bgrect y1 invalid: %v", err)
			}
			x2, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				return nil, fmt.Errorf("bgrect x2 invalid: %v", err)
			}
			y2, err := strconv.ParseFloat(args[3], 64)
			if err != nil {
				return nil, fmt.Errorf("bgrect y2 invalid: %v", err)
			}
			res = append(res, painter.BgrectOp{X1: x1, Y1: y1, X2: x2, Y2: y2})

		case "figure":
			if len(args) != 2 {
				return nil, fmt.Errorf("figure command requires 2 arguments, got %d", len(args))
			}
			x, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return nil, fmt.Errorf("figure x invalid: %v", err)
			}
			y, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return nil, fmt.Errorf("figure y invalid: %v", err)
			}
			res = append(res, painter.FigureOp{X: x, Y: y})

		case "move":
			if len(args) != 2 {
				return nil, fmt.Errorf("move command requires 2 arguments, got %d", len(args))
			}
			x, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return nil, fmt.Errorf("move x invalid: %v", err)
			}
			y, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return nil, fmt.Errorf("move y invalid: %v", err)
			}
			res = append(res, painter.MoveOp{X: x, Y: y})

		case "reset":
			if len(args) != 0 {
				return nil, fmt.Errorf("reset command takes no arguments, got %d", len(args))
			}
			res = append(res, painter.ResetOp)

		default:
			return nil, fmt.Errorf("unknown command: %s", instruction)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return res, nil
}