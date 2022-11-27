package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"platformer/common"
)

const playingState = "playing"
const dyingState = "dying"

const dyingTimeAmount = 0.8
const shootTimerAmount = 0.5
const moveFrameAmount = 0.04
const playerYNormal = 200.0

type Player struct {
	x           float64
	y           float64
	targetY     float64
	frame       int
	frames      int
	sizex       int
	sizey       int
	drawSizex   int
	drawSizey   int
	state       string
	speed       float64
	timer       float64
	shootTimer  float64
	image       *ebiten.Image
	lives       int
	animTimer   float64
	targetFrame int
	moveYSpeed  float64
}

func NewPlayer(game *Game) *Player {
	p := &Player{
		state:     playingState,
		y:         common.ScreenHeight / 2,
		x:         common.ScreenWidth / 2,
		speed:     80,
		sizex:     16, // physical size
		sizey:     24, // physical size
		drawSizex: 32, // just for drawing
		drawSizey: 32, // just for drawing
		lives:     2,
		image:     game.images["player"],
	}
	return p
}

func (r *Player) Update(delta float64, game *Game) {
	switch r.state {
	case playingState:
		inputX := 0.0
		inputY := 0.0
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			inputX = -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			inputX = 1
		}
		r.x = r.x + (inputX * delta * r.speed)

		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			inputY = 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			inputY = -1
		}
		r.y = r.y + (inputY * delta * r.speed)
	}
}

func (r *Player) Draw(camera common.Camera) {
	if r.state == dyingState {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.x, r.y)
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
