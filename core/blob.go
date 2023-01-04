package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
	"math/rand"
	"platformer/common"
)

const (
	thinkStateIdle   = "idle"
	thinkStateTarget = "target"
)

type BlobEnemy struct {
	x                float64
	y                float64
	sizeX            float64
	currentAnimation string
	animations       map[string]*Animation
	health           int
	// ai
	directionX       int
	targetX          float64
	targetY          float64
	moveSpeed        float64
	hurtTimer        float64
	hurtAmountTime   float64
	thinkState       string
	lastKnownPlayerX float64
	tryJumpTimer     float64
	velocityY        float64
	jumpTimer        float64
	touchingGround   bool
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
			"jump": {
				image:           game.res.GetImage("blob-jump"),
				numFrames:       1,
				size:            32,
				frameTimeAmount: 1,
				isLoop:          true,
			},
		},
		health:         2,
		directionX:     1,
		moveSpeed:      80,
		hurtAmountTime: 0.4,
		thinkState:     thinkStateIdle,
		tryJumpTimer:   rand.Float64() * 100,
	}
}

func (r *BlobEnemy) Update(delta float64, game *Game) {
	cb := r.GetCollisionBox()
	if common.Overlap(game.Player.x+4, game.Player.y+8, 8, 8, cb.x, cb.y, cb.w, cb.h) {
		game.Player.TakeDamage(game)
	}
	if r.hurtTimer > 0 {
		r.currentAnimation = "hurt"
		r.hurtTimer = r.hurtTimer - delta
	} else {
		r.think(game)
		r.move(delta, game)
	}
	r.animations[r.currentAnimation].Update(delta)
}

const blobJumpHeight = 2.0 * common.TileSize
const blobJumpTime = 0.5
const timeBetweenJumps = 2.0

func (r *BlobEnemy) move(delta float64, game *Game) {
	moveX := 0.0

	gravity := (blobJumpHeight * -2) / (blobJumpTime * blobJumpTime)

	moveY := (r.velocityY * delta) + (0.5 * gravity * delta * delta)
	r.velocityY = r.velocityY + (gravity * delta)

	actualSpeed := r.moveSpeed
	if r.touchingGround {
		actualSpeed = 20
	}
	if math.Abs(r.x-r.targetX) < (r.moveSpeed * delta) {
		r.x = r.targetX
		r.currentAnimation = "idle"
	}
	if r.x < r.targetX {
		moveX = actualSpeed * delta
		r.currentAnimation = "run"
	}
	if r.x > r.targetX {
		moveX = -(actualSpeed * delta)
		r.currentAnimation = "run"
	}

	// alter movement after checking or collision
	cb := r.GetCollisionBox()
	oldX := cb.x
	newX := cb.x + moveX
	oldY := cb.y + cb.h
	newY := cb.y + cb.h - moveY

	// x movement collision
	tx, ty := int((newX)/common.TileSize), int(oldY/common.TileSize)
	game.debug.DrawBox(color.RGBA{R: 244, G: 12, B: 9, A: 244}, float64(tx*common.TileSize), float64(ty*common.TileSize), common.TileSize, common.TileSize)
	td := game.Level.tiledGrid.GetTileData(tx, ty)
	if td.Block || td.Damage || td.Platform {
		newX = oldX
	}
	tx, ty = int((newX+cb.w)/common.TileSize), int(oldY/common.TileSize)
	game.debug.DrawBox(color.RGBA{R: 244, G: 12, B: 9, A: 244}, float64(tx*common.TileSize), float64(ty*common.TileSize), common.TileSize, common.TileSize)
	td = game.Level.tiledGrid.GetTileData(tx, ty)
	if td.Block || td.Damage || td.Platform {
		newX = oldX
	}

	// y movement collision
	r.touchingGround = false
	tx, ty = int((oldX)/common.TileSize), int(newY/common.TileSize)
	game.debug.DrawBox(color.RGBA{R: 244, G: 12, B: 9, A: 244}, float64(tx*common.TileSize), float64(ty*common.TileSize), common.TileSize, common.TileSize)
	td = game.Level.tiledGrid.GetTileData(tx, ty)
	if td.Block || td.Damage || td.Platform {
		newY = oldY
		r.velocityY = 0
		r.touchingGround = true
	}
	tx, ty = int((oldX+cb.w)/common.TileSize), int(newY/common.TileSize)
	game.debug.DrawBox(color.RGBA{R: 244, G: 12, B: 9, A: 244}, float64(tx*common.TileSize), float64(ty*common.TileSize), common.TileSize, common.TileSize)
	td = game.Level.tiledGrid.GetTileData(tx, ty)
	if td.Block || td.Damage || td.Platform {
		newY = oldY
		r.velocityY = 0
		r.touchingGround = true
	}

	// jumping
	r.tryJumpTimer = r.tryJumpTimer + delta
	if r.touchingGround && r.tryJumpTimer > timeBetweenJumps {
		r.tryJumpTimer = 0
		r.velocityY = (2 * blobJumpHeight) / blobJumpTime
	}
	if r.velocityY > 0 {
		r.currentAnimation = "jump"
	}

	r.x = newX - 8
	r.y = newY - 8 - cb.h
}

const (
	blobViewDistance = common.TileSize * 6
)

func (r *BlobEnemy) think(game *Game) {
	canSeePlayer := false
	if r.x < game.Player.x+blobViewDistance && r.x > game.Player.x-blobViewDistance {
		if r.y < game.Player.y+blobViewDistance && r.y > game.Player.y-blobViewDistance {
			canSeePlayer = true
			r.lastKnownPlayerX = game.Player.x
		}
	}
	atTarget := false
	if r.targetX == r.x {
		atTarget = true
	}
	switch r.thinkState {
	case thinkStateIdle:
		r.targetX = r.x
		r.targetY = r.y
		if canSeePlayer {
			r.thinkState = thinkStateTarget
		}
	case thinkStateTarget:
		if r.touchingGround {
			r.targetX = r.lastKnownPlayerX
		}
		if !canSeePlayer && atTarget {
			r.thinkState = thinkStateIdle
		}
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
