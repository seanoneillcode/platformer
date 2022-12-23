package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
	"platformer/common"
)

const hurtAmountTime = 0.4

type Crawler struct {
	x                float64
	y                float64
	sizeX            float64
	currentAnimation string
	animations       map[string]*Animation
	health           int
	// ai
	directionX int
	targetX    float64
	targetY    float64
	moveSpeed  float64
	hurtTimer  float64
}

func NewCrawler(x float64, y float64, game *Game) *Crawler {
	return &Crawler{
		x:                x,
		y:                y,
		sizeX:            24,
		currentAnimation: "run",
		animations: map[string]*Animation{
			"run": {
				image:           game.images["crawler-run"],
				numFrames:       5,
				size:            24,
				frameTimeAmount: 0.1,
				isLoop:          true,
			},
			"idle": {
				image:           game.images["crawler-idle"],
				numFrames:       2,
				size:            24,
				frameTimeAmount: 0.4,
				isLoop:          true,
			},
			"hurt": {
				image:           game.images["crawler-hurt"],
				numFrames:       2,
				size:            24,
				frameTimeAmount: 0.1,
				isLoop:          true,
			},
			//"death": {
			//	image:           game.images["crawler-death"],
			//	numFrames:       2,
			//	size:            24,
			//	frameTimeAmount: 0.4,
			//	isLoop:          false,
			//},
		},
		health:     3,
		directionX: 1,
		moveSpeed:  16,
	}
}

func (r *Crawler) Update(delta float64, game *Game) {
	if common.Overlap(game.player.x+4, game.player.y+8, 8, 8, r.x+2, r.y+2, 12, 12) {
		game.player.TakeDamage(game)
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

func (r *Crawler) think(game *Game) {
	if r.directionX > 0 {
		tx, ty := int(r.x/common.TileSize)+1, int(r.y/common.TileSize)
		game.debug.DrawBox(color.RGBA{R: 244, G: 12, B: 9, A: 244}, float64(tx*common.TileSize), float64(ty*common.TileSize), common.TileSize, common.TileSize)

		td := game.level.tiledGrid.GetTileData(tx, ty)
		if td.Block || td.Damage || td.Platform {
			r.directionX = r.directionX * -1
			return
		}

		// check tile below
		tx, ty = int((r.x)/common.TileSize)+1, int((r.y/common.TileSize)+1)
		td = game.level.tiledGrid.GetTileData(tx, ty)
		game.debug.DrawBox(color.RGBA{R: 120, G: 12, B: 44, A: 244}, float64(tx*common.TileSize), float64(ty*common.TileSize), common.TileSize, common.TileSize)
		if td.Block || td.Platform {
			r.targetX = float64(tx*common.TileSize) + float64(r.directionX*common.TileSize)
			return
		}

		r.directionX = r.directionX * -1
	} else {
		tx, ty := int(r.x/common.TileSize), int(r.y/common.TileSize)
		game.debug.DrawBox(color.RGBA{R: 244, G: 12, B: 9, A: 244}, float64(tx*common.TileSize), float64(ty*common.TileSize), common.TileSize, common.TileSize)

		td := game.level.tiledGrid.GetTileData(tx, ty)
		if td.Block || td.Damage || td.Platform {
			r.directionX = r.directionX * -1
			return
		}

		// check tile below
		tx, ty = int((r.x)/common.TileSize), int((r.y/common.TileSize)+1)
		td = game.level.tiledGrid.GetTileData(tx, ty)
		game.debug.DrawBox(color.RGBA{R: 120, G: 12, B: 44, A: 244}, float64(tx*common.TileSize), float64(ty*common.TileSize), common.TileSize, common.TileSize)
		if td.Block || td.Platform {
			r.targetX = float64(tx*common.TileSize) + float64(r.directionX*common.TileSize)
			return
		}

		r.directionX = r.directionX * -1
	}
}

func (r *Crawler) Draw(camera common.Camera) {

	op := &ebiten.DrawImageOptions{}
	if r.directionX > 0 {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(r.sizeX, 0)
	}
	op.GeoM.Translate(r.x-4, r.y-8)
	op.GeoM.Scale(common.Scale, common.Scale)
	camera.DrawImage(r.animations[r.currentAnimation].GetCurrentFrame(), op)
}

func (r *Crawler) GetHurt(game *Game) {
	r.health = r.health - 1
	r.hurtTimer = hurtAmountTime
	r.animations["hurt"].Play()
	if r.health == 0 {
		game.level.RemoveEnemy(r)
		// play effect
	}
}

func (r *Crawler) GetCollisionBox() CollisionBox {
	return CollisionBox{
		x: r.x + 2,
		y: r.y + 2,
		w: 12,
		h: 12,
	}
}
