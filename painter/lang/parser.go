package lang

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/DmytroHalai/kpi-3/painter"
)

type Parser struct {
}

func (p *Parser) Parse(in io.Reader, scene *painter.Scene) ([]painter.Operation, error) {
	var res []painter.Operation
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		cmd, args := parts[0], parts[1:]

		op, err := p.parseCommand(cmd, args, scene)
		if err != nil {
			return nil, err
		}
		res = append(res, op...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (p *Parser) parseCommand(cmd string, args []string, scene *painter.Scene) ([]painter.Operation, error) {
	switch cmd {
	case "white", "green", "update", "reset":
		if len(args) != 0 {
			return nil, fmt.Errorf("%s command takes no arguments, got %d", cmd, len(args))
		}
		return []painter.Operation{p.simpleOp(cmd, scene)}, nil

	case "bgrect":
		if len(args) != 4 {
			return nil, fmt.Errorf("bgrect command requires 4 arguments, got %d", len(args))
		}
		ints, err := parseArgsAsInts(args)
		if err != nil {
			return nil, fmt.Errorf("bgrect arg error: %v", err)
		}
		return []painter.Operation{painter.BgRectOp(scene, ints[0], ints[1], ints[2], ints[3])}, nil

	case "figure":
		if len(args) != 2 {
			return nil, fmt.Errorf("figure command requires 2 arguments, got %d", len(args))
		}
		ints, err := parseArgsAsInts(args)
		if err != nil {
			return nil, fmt.Errorf("figure arg error: %v", err)
		}
		return []painter.Operation{painter.ShapeOp(scene, ints[0], ints[1])}, nil

	case "move":
		if len(args) != 2 {
			return nil, fmt.Errorf("move command requires 2 arguments, got %d", len(args))
		}
		ints, err := parseArgsAsInts(args)
		if err != nil {
			return nil, fmt.Errorf("move arg error: %v", err)
		}
		return []painter.Operation{painter.MoveOp(scene, ints[0], ints[1])}, nil

	default:
		return nil, fmt.Errorf("unknown command: %s", cmd)
	}
}

func (p *Parser) simpleOp(cmd string, scene *painter.Scene) painter.Operation {
	switch cmd {
	case "white":
		return painter.WhiteFill(scene)
	case "green":
		return painter.GreenFill(scene)
	case "update":
		return painter.UpdateOp
	case "reset":
		return painter.ResetOp(scene)
	default:
		panic("unreachable")
	}
}

func parseArgsAsInts(args []string) ([]int, error) {
	ints := make([]int, len(args))
	for i, arg := range args {
		val, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return nil, fmt.Errorf("arg %d invalid: %v", i+1, err)
		}
		ints[i] = int(val)
	}
	return ints, nil
}
