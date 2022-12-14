package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"platformer/common"
)

type EffectSprite struct {
	x           float64
	y           float64
	w           float64
	h           float64
	rot         float64
	animation   *Animation
	ttl         float64
	isTemporary bool
	isFlipX     bool
}

func (r *EffectSprite) Update(delta float64, game *Game) {
	r.animation.Update(delta)
	if r.isTemporary {
		r.ttl = r.ttl - delta
		if r.ttl < 0 {
			game.RemoveEffectSprite(r)
		}
	}
}

func (r *EffectSprite) Draw(camera common.Camera) {
	op := &ebiten.DrawImageOptions{}
	if r.rot != 0 {
		offset := r.w / 2
		op.GeoM.Translate(-offset, -offset)
		op.GeoM.Rotate(r.rot)
		op.GeoM.Translate(offset, offset)
	}
	if r.isFlipX {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(r.w, 0)
	}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	camera.DrawImage(r.animation.GetCurrentFrame(), op)
}
