package painter_test

import (
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/roman-mazur/architecture-lab-3/painter/lang"
	"github.com/roman-mazur/architecture-lab-3/ui"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/paint"
)

type testReceiver struct {
	texture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.texture = t
}

func (tr *testReceiver) Execute(v *ui.Visualizer) {}

type MockScreen struct{}

func (ms *MockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	return &MockWindow{}, nil
}

func (ms *MockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	return nil, nil
}

func (ms *MockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return &MockTexture{size: size}, nil
}

type MockWindow struct{}

func (mw *MockWindow) Release()                                                                     {}
func (mw *MockWindow) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle)                 {}
func (mw *MockWindow) Fill(dr image.Rectangle, src color.Color, op draw.Op)                         {}
func (mw *MockWindow) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op) {}
func (mw *MockWindow) Draw(aff3 [6]float64, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions) {
}
func (mw *MockWindow) DrawUniform(aff3 [6]float64, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions) {
}
func (mw *MockWindow) Publish()               {}
func (mw *MockWindow) Send(event interface{}) {}
func (mw *MockWindow) NextEvent() interface{} { return nil }
func (mw *MockWindow) Copy(dp image.Point, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions) {
}

type MockTexture struct {
	size image.Point
}

func (mt *MockTexture) Release()                                                     {}
func (mt *MockTexture) Size() image.Point                                            { return mt.size }
func (mt *MockTexture) Bounds() image.Rectangle                                      { return image.Rect(0, 0, mt.size.X, mt.size.Y) }
func (mt *MockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}
func (mt *MockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op)         {}
func (mt *MockTexture) Scale(dr image.Rectangle, src screen.Texture, sr image.Rectangle, op draw.Op) {
}
func (mt *MockTexture) Draw(aff3 [6]float64, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions) {
}
func (mt *MockTexture) DrawUniform(aff3 [6]float64, src screen.Texture, sr image.Rectangle, op draw.Op, opts *screen.DrawOptions) {
}

var WhiteFill = color.RGBA{255, 255, 255, 255}
var GreenFill = color.RGBA{0, 255, 0, 255}

type UpdateOp struct{}

func (op UpdateOp) Execute(v *ui.Visualizer) {
	v.W.Send(paint.Event{})
}

func TestLoop(t *testing.T) {
	var tr testReceiver
	var l painter.Loop
	l.Receiver = &tr

	mockScreen := &MockScreen{}
	l.Start(mockScreen)

	l.Post(painter.BgRectOp{Rect: image.Rect(0, 0, 100, 100), FillColor: WhiteFill})
	l.Post(painter.FigureOp{Rect: image.Rect(10, 10, 90, 90), FillColor: GreenFill})
	l.Post(UpdateOp{})

	l.StopAndWait()

	if tr.texture == nil {
		t.Fatal("expected texture to be updated")
	}
}

func TestParser(t *testing.T) {
	p := &lang.Parser{}
	cmds, err := p.Parse(strings.NewReader("BGRECT 0 0 100 100 ffffff\nFIGURE 10 10 90 90 00ff00\nMOVE 10 10"))
	assert.Nil(t, err)
	assert.Len(t, cmds, 3)
}
