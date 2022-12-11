package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"platformer/common"
)

const playingState = "playing"
const dyingState = "dying"

const standardJumpHeight = 16 * 3.1
const standardJumpTime = 0.4
const standardFallTime = 0.36
const minimumJumpHeight = 16
const coyoteTimeAmount = 0.16
const fudge = 0.001
const runAcc = 20.0
const maxRunVelocity = 100
const tryJumpMarginTime = 0.1

type Player struct {
	x                  float64
	y                  float64
	targetY            float64
	sizex              float64
	sizey              float64
	drawOffsetX        float64
	drawOffsetY        float64
	drawSizex          int
	drawSizey          int
	state              string
	timer              float64
	shootTimer         float64
	image              *ebiten.Image
	lives              int
	animTimer          float64
	targetFrame        int
	moveYSpeed         float64
	jumpAcc            float64
	velocityY          float64
	velocityX          float64
	coyoteTimer        float64
	jumpTimer          float64
	wasPressingJump    bool
	alreadyAbortedJump bool
	direction          int
	targetVelocityX    float64
	isFlip             bool
	currentAnimation   string
	animations         map[string]*Animation
	tryJumpTimer       float64
}

func NewPlayer(game *Game) *Player {
	p := &Player{
		state:            playingState,
		y:                common.ScreenHeight / 2,
		x:                common.ScreenWidth / 2,
		sizex:            15, // physical size
		sizey:            16, // physical size
		drawOffsetX:      8,  // just for drawing
		drawOffsetY:      16, // just for drawing
		drawSizex:        32, // just for drawing
		drawSizey:        32,
		lives:            2,
		currentAnimation: "idle",
		animations: map[string]*Animation{
			"run": {
				image:           game.images["player-run"],
				numFrames:       6,
				size:            32,
				frameTimeAmount: 0.1,
				isLoop:          true,
			},
			"idle": {
				image:           game.images["player-idle"],
				numFrames:       1,
				size:            32,
				frameTimeAmount: 1,
				isLoop:          true,
			},
			"jump": {
				image:           game.images["player-jump"],
				numFrames:       1,
				size:            32,
				frameTimeAmount: 1,
				isLoop:          true,
			},
			"fall": {
				image:           game.images["player-fall"],
				numFrames:       1,
				size:            32,
				frameTimeAmount: 1,
				isLoop:          true,
			},
		},
		velocityY:       0,
		velocityX:       0,
		direction:       1,
		targetVelocityX: 0,
	}
	return p
}

func (r *Player) Update(delta float64, game *Game) {
	r.animations[r.currentAnimation].Update(delta)
	switch r.state {
	case playingState:
		var tryJump bool
		var pressJump bool
		var tryFall bool
		r.targetVelocityX = 0
		r.currentAnimation = "idle"
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			r.direction = -1
			r.targetVelocityX = -maxRunVelocity
			r.isFlip = true
			r.currentAnimation = "run"
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			r.direction = 1
			r.targetVelocityX = maxRunVelocity
			r.isFlip = false
			r.currentAnimation = "run"
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			tryFall = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {

		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			pressJump = true
		}
		r.tryJumpTimer = r.tryJumpTimer - delta
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			tryJump = true
			r.tryJumpTimer = tryJumpMarginTime
		}

		time := standardJumpTime
		if r.jumpTimer > (standardJumpTime) {
			time = standardFallTime
		}
		gravity := (standardJumpHeight * -2) / (time * time)

		oldx := r.x
		oldy := r.y

		newx := r.x + (delta * r.velocityX)
		movey := (r.velocityY * delta) + (0.5 * gravity * delta * delta)
		newy := r.y - movey
		r.velocityY = r.velocityY + (gravity * delta)

		var hitCeiling = false
		var hitFloor = false
		var hitWall = false
		td := game.level.tiledGrid.GetTileData(int(newx/common.TileSize), int(oldy/common.TileSize))
		if td.Block {
			newx = float64(td.X*common.TileSize) + common.TileSize + fudge
			hitWall = true
		}
		td = game.level.tiledGrid.GetTileData(int((newx+r.sizex)/common.TileSize), int(oldy/common.TileSize))
		if td.Block {
			newx = float64(td.X*common.TileSize) - r.sizex - fudge
			hitWall = true
		}
		td = game.level.tiledGrid.GetTileData(int(newx/common.TileSize), int((oldy+r.sizey)/common.TileSize))
		if td.Block {
			newx = float64(td.X*common.TileSize) + common.TileSize + fudge
			hitWall = true
		}
		td = game.level.tiledGrid.GetTileData(int((newx+r.sizex)/common.TileSize), int((oldy+r.sizey)/common.TileSize))
		if td.Block {
			newx = float64(td.X*common.TileSize) - r.sizex - fudge
			hitWall = true
		}

		td = game.level.tiledGrid.GetTileData(int(oldx/common.TileSize), int((newy-4)/common.TileSize))
		if td.Block {
			newy = float64(td.Y*common.TileSize) + common.TileSize + fudge + 4
			if newy > 0 {
				hitCeiling = true
			}
		}
		td = game.level.tiledGrid.GetTileData(int(oldx/common.TileSize), int((newy+r.sizey)/common.TileSize))
		if td.Block {
			distance := float64(td.Y*common.TileSize) - (oldy + r.sizey)
			newy = oldy + distance - fudge
			hitFloor = true
			r.velocityY = 0
			r.coyoteTimer = coyoteTimeAmount
		}
		if !tryFall && td.Platform && newy > oldy {
			distance := float64(td.Y*common.TileSize) - (oldy + r.sizey)
			if distance > -1 {
				newy = oldy + distance - fudge
				hitFloor = true
				r.velocityY = 0
				r.coyoteTimer = coyoteTimeAmount
			}
		}
		td = game.level.tiledGrid.GetTileData(int((oldx+r.sizex)/common.TileSize), int((newy-4)/common.TileSize))
		if td.Block {
			newy = float64(td.Y*common.TileSize) + common.TileSize + fudge + 4
			if newy > 0 {
				hitCeiling = true
			}
		}
		td = game.level.tiledGrid.GetTileData(int((oldx+r.sizex)/common.TileSize), int((newy+r.sizey)/common.TileSize))
		if td.Block {
			distance := float64(td.Y*common.TileSize) - (oldy + r.sizey)
			newy = oldy + distance - fudge
			hitFloor = true
			r.velocityY = 0
			r.coyoteTimer = coyoteTimeAmount
		}
		if !tryFall && td.Platform && newy > oldy {
			distance := float64(td.Y*common.TileSize) - (oldy + r.sizey)
			if distance > -1 {
				newy = oldy + distance - fudge
				hitFloor = true
				r.velocityY = 0
				r.coyoteTimer = coyoteTimeAmount
			}
		}

		r.x = newx
		r.y = newy

		if tryJump || r.tryJumpTimer > 0 {
			if hitFloor || r.coyoteTimer > 0 {
				r.coyoteTimer = 0
				r.velocityY = (2 * standardJumpHeight) / standardJumpTime
				r.jumpTimer = 0
			}
		}
		if !pressJump {
			// if player is currently jumping in the first half phase of jumping
			if r.jumpTimer < standardJumpTime && r.wasPressingJump && !r.alreadyAbortedJump {
				r.alreadyAbortedJump = true
				r.velocityY = (2 * minimumJumpHeight) / (standardJumpTime)
			}
		}
		r.jumpTimer = r.jumpTimer + delta
		if r.coyoteTimer > 0 {
			r.coyoteTimer = r.coyoteTimer - delta
		}
		if hitFloor {
			r.alreadyAbortedJump = false
		}
		if oldy != newy {
			if r.velocityY < 0 {
				r.currentAnimation = "fall"
			}
			if r.velocityY > 0 {
				r.currentAnimation = "jump"
			}
		}
		if hitCeiling {
			r.alreadyAbortedJump = true
			r.velocityY = -20
		}
		r.wasPressingJump = pressJump
		if hitWall {
			r.targetVelocityX = 0
			r.velocityX = 0

		}

		if r.velocityX < r.targetVelocityX {
			r.velocityX = r.velocityX + runAcc
			if r.velocityX > r.targetVelocityX {
				r.velocityX = r.targetVelocityX
			}
		}
		if r.velocityX > r.targetVelocityX {
			r.velocityX = r.velocityX - runAcc
			if r.velocityX < r.targetVelocityX {
				r.velocityX = r.targetVelocityX
			}
		}
	}
}

func (r *Player) Draw(camera common.Camera) {
	if r.state == dyingState {
		return
	}
	op := &ebiten.DrawImageOptions{}
	if r.isFlip {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(r.drawSizex), 0)
	}

	op.GeoM.Translate(r.x-r.drawOffsetX, r.y-r.drawOffsetY)
	op.GeoM.Scale(common.Scale, common.Scale)

	camera.DrawImage(r.animations[r.currentAnimation].GetCurrentFrame(), op)
}

func (r *Player) GetHit(game *Game) {
	r.lives = r.lives - 1
	if r.lives < 0 {
		r.state = dyingState
	}
}

func (r *Player) GetPos() (float64, float64) {
	return r.x, r.y
}
