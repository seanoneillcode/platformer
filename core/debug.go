package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"platformer/common"
)

type DebugDrawer struct {
	boxes     []*debugBox
	image     *ebiten.Image
	showDebug bool
}

func NewDebug() *DebugDrawer {
	return &DebugDrawer{
		boxes:     []*debugBox{},
		showDebug: false,
		image:     common.LoadImage("debug-pixel.png"),
	}
}

type debugBox struct {
	x float64
	y float64
	w float64
	h float64
	c color.Color
}

func (r *DebugDrawer) Update(delta float64, game *Game) {
	if len(r.boxes) > 0 {
		r.boxes = []*debugBox{}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		r.showDebug = !r.showDebug
	}
}

func (r *DebugDrawer) Draw(camera common.Camera) {
	if !r.showDebug {
		return
	}
	for _, box := range r.boxes {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(box.x, box.y)
		op.GeoM.Scale(common.Scale, common.Scale)
		op.GeoM.Scale(box.w/16.0, box.h/16.0)
		op.ColorM.ScaleWithColor(box.c)
		camera.DrawImage(r.image, op)
	}
}

func (r *DebugDrawer) DrawBox(c color.Color, x, y, w, h float64) {
	r.boxes = append(r.boxes, &debugBox{
		x: x,
		y: y,
		w: w,
		h: h,
		c: c,
	})
}
