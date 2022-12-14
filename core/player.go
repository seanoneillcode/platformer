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
const standardFallTime = 0.4
const minimumJumpHeight = 16
const coyoteTimeAmount = 0.16
const fudge = 0.001
const runAcc = 20.0
const maxRunVelocity = 100
const ladderVelocity = 70
const tryJumpMarginTime = 0.1
const ladderGrabAllowance = 8.0

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
	targetVelocityX    float64
	isFlip             bool
	currentAnimation   string
	animations         map[string]*Animation
	tryJumpTimer       float64
	lockedToLadder     bool
}

func NewPlayer(game *Game) *Player {
	p := &Player{
		state:            playingState,
		x:                19 * common.TileSize,
		y:                12 * common.TileSize,
		sizex:            16, // physical size
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
			"climb": {
				image:           game.images["player-climb"],
				numFrames:       1,
				size:            32,
				frameTimeAmount: 1,
				isLoop:          true,
			},
		},
		velocityY:       0,
		velocityX:       0,
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
		var tryMovey = 0.0
		r.targetVelocityX = 0
		r.currentAnimation = "idle"

		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			r.targetVelocityX = -maxRunVelocity
			r.isFlip = true
			r.currentAnimation = "run"
			if r.lockedToLadder {
				r.targetVelocityX = -maxRunVelocity / 2.0
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			r.targetVelocityX = maxRunVelocity
			r.isFlip = false
			r.currentAnimation = "run"
			if r.lockedToLadder {
				r.targetVelocityX = maxRunVelocity / 2.0
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			tryFall = true
			tryMovey = 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			tryMovey = -1
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

		partial := 2.0

		newx := r.x + (delta * r.velocityX)
		movey := 0.0
		if !r.lockedToLadder {
			movey = (r.velocityY * delta) + (0.5 * gravity * delta * delta)
		}
		newy := r.y - movey
		r.velocityY = r.velocityY + (gravity * delta)

		var hitWall = false
		tx, ty := int((newx+partial)/common.TileSize), int(oldy/common.TileSize)
		td := game.level.tiledGrid.GetTileData(tx, ty)
		if td.Block {
			newx = float64(tx*common.TileSize) + common.TileSize + fudge - partial
			hitWall = true
		}

		tx, ty = int(((newx+partial)+(r.sizex-partial))/common.TileSize), int(oldy/common.TileSize)
		td = game.level.tiledGrid.GetTileData(tx, ty)
		if td.Block {
			newx = float64(tx*common.TileSize) - (r.sizex - partial) - fudge - partial
			hitWall = true
		}

		tx, ty = int((newx+partial)/common.TileSize), int((oldy+r.sizey)/common.TileSize)
		td = game.level.tiledGrid.GetTileData(tx, ty)
		if td.Block {
			newx = float64(tx*common.TileSize) + common.TileSize + fudge - partial
			hitWall = true
		}

		tx, ty = int(((newx+partial)+(r.sizex-partial))/common.TileSize), int((oldy+r.sizey)/common.TileSize)
		td = game.level.tiledGrid.GetTileData(tx, ty)
		if td.Block {
			newx = float64(tx*common.TileSize) - (r.sizex - partial) - fudge - partial
			hitWall = true
		}

		var hitCeiling = false
		var hitFloor = false

		tx, ty = int((oldx+partial)/common.TileSize), int((newy+r.sizey)/common.TileSize)
		td = game.level.tiledGrid.GetTileData(tx, ty)
		if td.Block {
			distance := float64(ty*common.TileSize) - (oldy + r.sizey)
			newy = oldy + distance - fudge
			hitFloor = true
			r.velocityY = 0
			r.coyoteTimer = coyoteTimeAmount
		}
		if !tryFall && td.Platform && newy > oldy {
			distance := float64(ty*common.TileSize) - (oldy + r.sizey)
			if distance > -1 {
				newy = oldy + distance - fudge
				hitFloor = true
				r.velocityY = 0
				r.coyoteTimer = coyoteTimeAmount
			}
		}

		tx, ty = int((oldx+partial)/common.TileSize), int((newy-4)/common.TileSize)
		td = game.level.tiledGrid.GetTileData(tx, ty)
		if td.Block {
			newy = float64(ty*common.TileSize) + common.TileSize + fudge + 4
			if newy > 0 {
				hitCeiling = true
			}
		}

		tx, ty = int(((oldx+partial)+(r.sizex-partial))/common.TileSize), int((newy+r.sizey)/common.TileSize)
		td = game.level.tiledGrid.GetTileData(tx, ty)
		if td.Block {
			distance := float64(ty*common.TileSize) - (oldy + r.sizey)
			newy = oldy + distance - fudge
			hitFloor = true
			r.velocityY = 0
			r.coyoteTimer = coyoteTimeAmount
		}
		if !tryFall && td.Platform && newy > oldy {
			distance := float64(ty*common.TileSize) - (oldy + r.sizey)
			if distance > -1 {
				newy = oldy + distance - fudge
				hitFloor = true
				r.velocityY = 0
				r.coyoteTimer = coyoteTimeAmount
			}
		}

		tx, ty = int(((oldx+partial)+(r.sizex-partial))/common.TileSize), int((newy-4)/common.TileSize)
		td = game.level.tiledGrid.GetTileData(tx, ty)
		if td.Block {
			newy = float64(ty*common.TileSize) + common.TileSize + fudge + 4
			if newy > 0 {
				hitCeiling = true
			}
		}

		var touchingLadder = false
		tx, ty = int((oldx+(r.sizex/2.0))/common.TileSize), int((oldy+4)/common.TileSize)
		td = game.level.tiledGrid.GetTileData(tx, ty)
		if td.Ladder {
			touchingLadder = true
			if tryMovey != 0 {
				middle := float64(tx * common.TileSize)
				left, right := middle-ladderGrabAllowance, middle+ladderGrabAllowance
				if oldx > left && oldx < right {
					r.lockedToLadder = true
					newy = oldy + (delta * ladderVelocity * tryMovey)
				}
			}

		}
		tx, ty = int((oldx+(r.sizex/2.0))/common.TileSize), int((oldy+r.sizey)/common.TileSize)
		td = game.level.tiledGrid.GetTileData(tx, ty)
		if td.Ladder {
			touchingLadder = true
			if tryMovey != 0 {
				middle := float64(tx * common.TileSize)
				left, right := middle-ladderGrabAllowance, middle+ladderGrabAllowance
				if oldx > left && oldx < right {
					r.lockedToLadder = true
					newy = oldy + (delta * ladderVelocity * tryMovey)

					// if move down, check for block
					if tryMovey == 1 {
						tx, ty = int((oldx+(r.sizex/2.0))/common.TileSize), int((newy+r.sizey)/common.TileSize)
						td = game.level.tiledGrid.GetTileData(tx, ty)
						if td.Block {
							r.lockedToLadder = false
							newy = oldy
						}
					}
				}
			}
		}
		if !touchingLadder {
			r.lockedToLadder = false
		}

		r.x = newx
		r.y = newy

		if tryJump || r.tryJumpTimer > 0 {
			if hitFloor || r.coyoteTimer > 0 || r.lockedToLadder {
				r.lockedToLadder = false
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
			r.lockedToLadder = false
		}
		if r.lockedToLadder {
			r.currentAnimation = "climb"
			r.velocityY = 0
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
		if r.lockedToLadder {
			//r.targetVelocityX = 0
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
