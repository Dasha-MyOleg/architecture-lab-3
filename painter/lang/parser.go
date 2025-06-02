package lang

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"io"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

// Parser визначає синтаксичний аналізатор для простого сценарію малювання.
type Parser struct{}

func (p *Parser) Parse(r io.Reader) ([]painter.Operation, error) {
	var ops []painter.Operation
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "BGRECT":
			if len(parts) != 6 {
				return nil, fmt.Errorf("invalid BGRECT command: %s", line)
			}
			x1, y1, x2, y2, err := parseRect(parts[1], parts[2], parts[3], parts[4])
			if err != nil {
				return nil, fmt.Errorf("invalid BGRECT coordinates: %v", err)
			}
			color, err := parseColor(parts[5])
			if err != nil {
				return nil, fmt.Errorf("invalid color: %v", err)
			}
			ops = append(ops, painter.BgRectOp{Rect: image.Rect(x1, y1, x2, y2), FillColor: color})

		case "FIGURE":
			if len(parts) != 6 {
				return nil, fmt.Errorf("invalid FIGURE command: %s", line)
			}
			x1, y1, x2, y2, err := parseRect(parts[1], parts[2], parts[3], parts[4])
			if err != nil {
				return nil, fmt.Errorf("invalid FIGURE coordinates: %v", err)
			}
			color, err := parseColor(parts[5])
			if err != nil {
				return nil, fmt.Errorf("invalid color: %v", err)
			}
			ops = append(ops, painter.FigureOp{Rect: image.Rect(x1, y1, x2, y2), FillColor: color})

		case "MOVE":
			if len(parts) != 3 {
				return nil, fmt.Errorf("invalid MOVE command: %s", line)
			}
			dx, dy, err := parsePoint(parts[1], parts[2])
			if err != nil {
				return nil, fmt.Errorf("invalid MOVE offset: %v", err)
			}
			ops = append(ops, painter.MoveOp{Offset: image.Pt(dx, dy)})

		default:
			return nil, fmt.Errorf("unknown command: %s", parts[0])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ops, nil
}

func parseRect(sx1, sy1, sx2, sy2 string) (x1, y1, x2, y2 int, err error) {
	x1, err = parseCoord(sx1)
	if err != nil {
		return
	}
	y1, err = parseCoord(sy1)
	if err != nil {
		return
	}
	x2, err = parseCoord(sx2)
	if err != nil {
		return
	}
	y2, err = parseCoord(sy2)
	return
}

func parsePoint(sx, sy string) (x, y int, err error) {
	x, err = parseCoord(sx)
	if err != nil {
		return
	}
	y, err = parseCoord(sy)
	return
}

func parseCoord(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

func parseColor(s string) (color.Color, error) {
	var r, g, b uint8
	_, err := fmt.Sscanf(s, "%02x%02x%02x", &r, &g, &b)
	return color.RGBA{R: r, G: g, B: b, A: 255}, err
}
