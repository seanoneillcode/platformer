package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"platformer/common"
)

const spellBulletSpeed = 200.0

type SpellObject struct {
	x         float64
	y         float64
	width     float64
	animation *Animation
	moveX     float64
	moveY     float64
	ttl       float64
	isFlip    bool
}

func NewSpellObject(game *Game, x, y, moveX, moveY float64) *SpellObject {
	return &SpellObject{
		x:      x,
		y:      y,
		moveX:  moveX,
		moveY:  moveY,
		ttl:    10,
		width:  16,
		isFlip: moveX < 0,
		animation: &Animation{
			image:           game.images["spell-bullet"],
			numFrames:       4,
			size:            16,
			frameTimeAmount: 0.1,
			isLoop:          true,
		},
	}
}

func (r *SpellObject) Update(delta float64, game *Game) {
	r.animation.Update(delta)
	r.x = r.x + (r.moveX * delta)
	r.y = r.y + (r.moveY * delta)
	r.ttl = r.ttl - delta
	if r.ttl < 0 {
		game.RemoveSpellObject(r)
		game.SpawnEffect(effectSpellHit, r.x, r.y)
		return
	}
	// check for collision with level, enemies, player etc
	tx, ty := int((r.x+8)/common.TileSize), int((r.y+8)/common.TileSize)
	td := game.level.tiledGrid.GetTileData(tx, ty)
	if td.Block {
		game.RemoveSpellObject(r)
		game.SpawnEffect(effectSpellHit, r.x, r.y)
		return
	}

	for _, e := range game.level.enemies {
		cb := e.GetCollisionBox()
		if common.Overlap(r.x+6, r.y+6, 4, 4, cb.x, cb.y, cb.w, cb.h) {
			e.GetHurt(game)
			game.RemoveSpellObject(r)
			game.SpawnEffect(effectSpellHit, r.x, r.y)
			return
		}
	}
}

func (r *SpellObject) Draw(camera common.Camera) {
	op := &ebiten.DrawImageOptions{}
	if r.isFlip {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(r.width, 0)
	}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	camera.DrawImage(r.animation.GetCurrentFrame(), op)
}
