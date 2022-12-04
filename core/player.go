package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"platformer/common"
)

const playingState = "playing"
const dyingState = "dying"

const standardJumpHeight = 16 * 3
const standardJumpTime = 0.4
const standardFallTime = 0.34
const minimumJumpHeight = 16
const coyoteTimeAmount = 0.16
const fudge = 0.15
const runAcc = 20.0
const maxRunVelocity = 120

type Player struct {
	x                  float64
	y                  float64
	targetY            float64
	frame              int
	frames             int
	sizex              float64
	sizey              float64
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
}

func NewPlayer(game *Game) *Player {
	p := &Player{
		state:           playingState,
		y:               common.ScreenHeight / 2,
		x:               common.ScreenWidth / 2,
		sizex:           12, // physical size
		sizey:           16, // physical size
		drawSizex:       16, // just for drawing
		drawSizey:       16, // just for drawing
		lives:           2,
		image:           game.images["player"],
		velocityY:       0,
		velocityX:       0,
		direction:       1,
		targetVelocityX: 0,
	}
	return p
}

func (r *Player) Update(delta float64, game *Game) {
	switch r.state {
	case playingState:
		var tryJump bool
		r.targetVelocityX = 0
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			r.direction = -1
			r.targetVelocityX = -maxRunVelocity
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			r.direction = 1
			r.targetVelocityX = maxRunVelocity
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {

		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {

		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			tryJump = true
		}

		time := standardJumpTime
		if r.jumpTimer > (standardJumpTime) {
			time = standardFallTime
		}
		gravity := (standardJumpHeight * -2) / (time * time)

		oldx := r.x
		oldy := r.y

		newx := r.x + (delta * r.velocityX)
		newy := r.y - (r.velocityY * delta) + (0.5 * gravity * delta * delta)
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

		td = game.level.tiledGrid.GetTileData(int(oldx/common.TileSize), int(newy/common.TileSize))
		if td.Block {
			newy = float64(td.Y*common.TileSize) + common.TileSize + fudge
			if newy > 0 {
				hitCeiling = true
			}
		}
		td = game.level.tiledGrid.GetTileData(int(oldx/common.TileSize), int((newy+r.sizey)/common.TileSize))
		if td.Block {
			newy = float64(td.Y*common.TileSize) - common.TileSize - fudge
			if newy > 0 {
				hitFloor = true
				r.velocityY = 0
				r.coyoteTimer = coyoteTimeAmount
			}
		}
		td = game.level.tiledGrid.GetTileData(int((oldx+r.sizex)/common.TileSize), int(newy/common.TileSize))
		if td.Block {
			newy = float64(td.Y*common.TileSize) + common.TileSize + fudge
			if newy > 0 {
				hitCeiling = true
			}
		}
		td = game.level.tiledGrid.GetTileData(int((oldx+r.sizex)/common.TileSize), int((newy+r.sizey)/common.TileSize))
		if td.Block {
			newy = float64(td.Y*common.TileSize) - common.TileSize - fudge
			if newy > 0 {
				hitFloor = true
				r.velocityY = 0
				r.coyoteTimer = coyoteTimeAmount
			}
		}

		r.x = newx
		r.y = newy

		if tryJump {
			if hitFloor || r.coyoteTimer > 0 {
				r.coyoteTimer = 0
				r.velocityY = (2 * standardJumpHeight) / standardJumpTime
				r.jumpTimer = 0
			}
		} else {
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
		if hitCeiling {
			r.alreadyAbortedJump = true
			r.velocityY = -20
		}
		r.wasPressingJump = tryJump
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
	op.GeoM.Translate(r.x-2, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	camera.DrawImage(r.image.SubImage(image.Rect(r.frame*r.drawSizex, 0, (r.frame+1)*r.drawSizex, r.drawSizey)).(*ebiten.Image), op)
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
