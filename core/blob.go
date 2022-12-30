package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
	"platformer/common"
)

type BlobEnemy struct {
	x                float64
	y                float64
	sizeX            float64
	currentAnimation string
	animations       map[string]*Animation
	health           int
	// ai
	directionX     int
	targetX        float64
	targetY        float64
	moveSpeed      float64
	hurtTimer      float64
	hurtAmountTime float64
}

func NewBlobEnemy(x float64, y float64, game *Game) *BlobEnemy {
	return &BlobEnemy{
		x:                x,
		y:                y,
		sizeX:            32,
		currentAnimation: "run",
		animations: map[string]*Animation{
			"run": {
				image:           game.res.GetImage("blob-run"),
				numFrames:       2,
				size:            32,
				frameTimeAmount: 0.2,
				isLoop:          true,
			},
			"idle": {
				image:           game.res.GetImage("blob-idle"),
				numFrames:       1,
				size:            32,
				frameTimeAmount: 1,
				isLoop:          true,
			},
			"hurt": {
				image:           game.res.GetImage("blob-hurt"),
				numFrames:       2,
				size:            32,
				frameTimeAmount: 0.1,
				isLoop:          true,
			},
			"attack": {
				image:           game.res.GetImage("blob-attack"),
				numFrames:       2,
				size:            32,
				frameTimeAmount: 0.5,
				isLoop:          true,
			},
		},
		health:         2,
		directionX:     1,
		moveSpeed:      32,
		hurtAmountTime: 0.4,
	}
}

func (r *BlobEnemy) Update(delta float64, game *Game) {
	cb := r.GetCollisionBox()
	if common.Overlap(game.Player.x+4, game.Player.y+8, 8, 8, cb.x, cb.y, cb.w, cb.h) {
		game.Player.TakeDamage(game)
	}
	r.currentAnimation = "idle"
	if r.hurtTimer > 0 {
		r.currentAnimation = "hurt"
		r.hurtTimer = r.hurtTimer - delta
	} else {
		if math.Abs(r.x-r.targetX) < (r.moveSpeed * delta) {
			r.x = r.targetX
			r.currentAnimation = "run"
		}
		if r.x < r.targetX {
			r.x = r.x + (r.moveSpeed * delta)
			r.currentAnimation = "run"
		}
		if r.x > r.targetX {
			r.x = r.x - (r.moveSpeed * delta)
			r.currentAnimation = "run"
		}
	}
	r.animations[r.currentAnimation].Update(delta)
	// thinking
	r.think(game)
}

func (r *BlobEnemy) think(game *Game) {
	cb := r.GetCollisionBox()
	if r.directionX > 0 {

		tx, ty := int((cb.x+cb.w)/common.TileSize), int(r.y/common.TileSize)
		game.debug.DrawBox(color.RGBA{R: 244, G: 12, B: 9, A: 244}, float64(tx*common.TileSize), float64(ty*common.TileSize), common.TileSize, common.TileSize)

		td := game.Level.tiledGrid.GetTileData(tx, ty)
		if td.Block || td.Damage || td.Platform {
			r.directionX = r.directionX * -1
			return
		}

		tx, ty = int((cb.x+cb.w)/common.TileSize), int(r.y/common.TileSize+1)
		game.debug.DrawBox(color.RGBA{R: 244, G: 12, B: 9, A: 244}, float64(tx*common.TileSize), float64(ty*common.TileSize), common.TileSize, common.TileSize)

		td = game.Level.tiledGrid.GetTileData(tx, ty)
		if td.Block || td.Damage || td.Platform {
			r.directionX = r.directionX * -1
			return
		}

		// check tile below
		tx, ty = int((cb.x+cb.w)/common.TileSize), int((r.y/common.TileSize)+2)
		td = game.Level.tiledGrid.GetTileData(tx, ty)
		game.debug.DrawBox(color.RGBA{R: 120, G: 12, B: 44, A: 244}, float64(tx*common.TileSize), float64(ty*common.TileSize), common.TileSize, common.TileSize)
		if td.Block || td.Platform {
			r.targetX = float64(tx*common.TileSize) + float64(r.directionX*common.TileSize)
			return
		}

		r.directionX = r.directionX * -1
	} else {
		tx, ty := int(cb.x/common.TileSize), int(r.y/common.TileSize)
		game.debug.DrawBox(color.RGBA{R: 244, G: 12, B: 9, A: 244}, float64(tx*common.TileSize), float64(ty*common.TileSize), common.TileSize, common.TileSize)

		td := game.Level.tiledGrid.GetTileData(tx, ty)
		if td.Block || td.Damage || td.Platform {
			r.directionX = r.directionX * -1
			return
		}

		tx, ty = int(cb.x/common.TileSize), int(r.y/common.TileSize)+1
		game.debug.DrawBox(color.RGBA{R: 244, G: 12, B: 9, A: 244}, float64(tx*common.TileSize), float64(ty*common.TileSize), common.TileSize, common.TileSize)

		td = game.Level.tiledGrid.GetTileData(tx, ty)
		if td.Block || td.Damage || td.Platform {
			r.directionX = r.directionX * -1
			return
		}

		// check tile below
		tx, ty = int((cb.x)/common.TileSize), int((r.y/common.TileSize)+2)
		td = game.Level.tiledGrid.GetTileData(tx, ty)
		game.debug.DrawBox(color.RGBA{R: 120, G: 12, B: 44, A: 244}, float64(tx*common.TileSize), float64(ty*common.TileSize), common.TileSize, common.TileSize)
		if td.Block || td.Platform {
			r.targetX = float64(tx*common.TileSize) + float64(r.directionX*common.TileSize)
			return
		}

		r.directionX = r.directionX * -1
	}
}

func (r *BlobEnemy) Draw(camera common.Camera) {

	op := &ebiten.DrawImageOptions{}
	if r.directionX > 0 {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(r.sizeX, 0)
	}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	camera.DrawImage(r.animations[r.currentAnimation].GetCurrentFrame(), op)
}

func (r *BlobEnemy) GetHurt(game *Game) {
	r.health = r.health - 1
	r.hurtTimer = r.hurtAmountTime
	r.animations["hurt"].Play()
	if r.health == 0 {
		game.SpawnEffect(effectBlobDeath, r.x, r.y, r.directionX > 0, 0)
		game.Level.RemoveEnemy(r)
		// play effect
	}
}

func (r *BlobEnemy) GetCollisionBox() CollisionBox {
	return CollisionBox{
		x: r.x + 8,
		y: r.y + 8,
		w: 16,
		h: 24,
	}
}
