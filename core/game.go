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
	timer            float64
	state            string
	images           map[string]*ebiten.Image
	level            int
	score            int
}

func NewGame() *Game {
	r := &Game{
		images: map[string]*ebiten.Image{
			"player": common.LoadImage("player.png"),
		},
		lastUpdateCalled: time.Now(),
	}
	r.StartNewGame()
	return r
}

func (r *Game) Update() error {
	delta := float64(time.Now().Sub(r.lastUpdateCalled).Milliseconds()) / 1000
	r.lastUpdateCalled = time.Now()

	r.player.Update(delta, r)

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return common.NormalEscapeError
	}

	return nil
}

func (r *Game) Draw(screen *ebiten.Image) {
	r.player.Draw(screen)
}

func (r *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.ScreenWidth * common.Scale, common.ScreenHeight * common.Scale
}
