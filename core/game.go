package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"platformer/common"
	"time"
)

type Game struct {
	lastUpdateCalled time.Time
	player           *Player
	images           map[string]*ebiten.Image
	camera           *Camera
	level            *Level
}

func NewGame() *Game {
	r := &Game{
		images: map[string]*ebiten.Image{
			//"player": common.LoadImage("player.png"),
			"player": common.LoadImage("test-player.png"),
		},
		lastUpdateCalled: time.Now(),
	}
	r.LoadLevel("test-level")
	return r
}

func (r *Game) Update() error {
	delta := float64(time.Now().Sub(r.lastUpdateCalled).Milliseconds()) / 1000
	r.lastUpdateCalled = time.Now()

	r.player.Update(delta, r)
	r.camera.Update(delta)
	r.level.Update(delta, r)

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return common.NormalEscapeError
	}

	return nil
}

func (r *Game) Draw(screen *ebiten.Image) {
	r.level.Draw(r.camera)
	r.player.Draw(r.camera)
	r.camera.DrawBuffer(screen)
}

func (r *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.ScreenWidth * common.Scale, common.ScreenHeight * common.Scale
}
