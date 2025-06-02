package painter

import (
	"image"
	"image/color"

	"github.com/roman-mazur/architecture-lab-3/ui"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/paint"
)

type Operation interface {
	Execute(v *ui.Visualizer)
}

type BgRectOp struct {
	FillColor color.Color
	Rect      image.Rectangle
}

func (op BgRectOp) Execute(v *ui.Visualizer) {
	v.BgColor = op.FillColor
	v.BgRect = op.Rect
	v.W.Send(paint.Event{})
}

type FigureOp struct {
	FillColor color.Color
	Rect      image.Rectangle
}

func (op FigureOp) Execute(v *ui.Visualizer) {
	v.W.Fill(op.Rect, op.FillColor, draw.Src)
	v.W.Send(paint.Event{})
}

type MoveOp struct {
	Offset image.Point
}

func (op MoveOp) Execute(v *ui.Visualizer) {
	v.Pos = v.Pos.Add(op.Offset)
	v.W.Send(paint.Event{})
}
