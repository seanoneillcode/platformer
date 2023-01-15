package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"platformer/common"
	"platformer/core"
)

const fadeOutTimeAmount = 1.0

type Popup struct {
	x         float64
	y         float64
	pageImage *ebiten.Image
	text      string
	//
	visible bool
	timer   float64
}

func (r *Popup) Update(delta float64, game *core.Game) {
	if r.timer > 0 {
		r.timer = r.timer - delta
		if r.timer < 0 {
			r.visible = false
		}
	}
}

func (r *Popup) Open(x float64, y float64, text string) {
	r.text = text
	r.x = x
	r.y = y
	r.visible = true
}

func (r *Popup) Close() {
	r.timer = fadeOutTimeAmount
}

func (r *Popup) Draw(screen *ebiten.Image) {
	if !r.visible {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(r.pageImage, op)

	paragraphX := r.x + 5
	paragraphY := r.y + 5
	common.DrawText(screen, r.text, paragraphX, paragraphY)
}
