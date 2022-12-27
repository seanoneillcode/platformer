package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"platformer/common"
)

type Flimsy struct {
	x     float64
	y     float64
	w     float64
	h     float64
	image *ebiten.Image
}

func (r *Flimsy) GetCollisionBox() CollisionBox {
	return CollisionBox{
		x: r.x,
		y: r.y,
		w: r.w,
		h: r.h,
	}
}

func (r *Flimsy) Draw(camera common.Camera) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	camera.DrawImage(r.image, op)
}
