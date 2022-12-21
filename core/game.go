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
	debug            *DebugDrawer
}

func NewGame() *Game {
	r := &Game{
		images: map[string]*ebiten.Image{
			"player-run":    common.LoadImage("player-run.png"),
			"player-idle":   common.LoadImage("player-idle.png"),
			"player-jump":   common.LoadImage("player-jump.png"),
			"player-fall":   common.LoadImage("player-fall.png"),
			"player-hurt":   common.LoadImage("player-hurt.png"),
			"player-death":  common.LoadImage("player-death.png"),
			"player-climb":  common.LoadImage("player-climb.png"),
			"book-pickup":   common.LoadImage("book.png"),
			"health-pickup": common.LoadImage("health.png"),
			"crawler-run":   common.LoadImage("crawler-run.png"),
			"crawler-idle":  common.LoadImage("crawler-idle.png"),
			"crawler-hurt":  common.LoadImage("crawler-hurt.png"),
			"crawler-die":   common.LoadImage("crawler-die.png"),
		},
		lastUpdateCalled: time.Now(),
		debug:            NewDebug(),
	}
	r.LoadLevel("long-level")
	return r
}

func (r *Game) Update() error {
	delta := float64(time.Now().Sub(r.lastUpdateCalled).Milliseconds()) / 1000
	r.lastUpdateCalled = time.Now()

	r.debug.Update(delta, r)
	r.player.Update(delta, r)
	r.level.Update(delta, r)
	r.camera.Update(delta, r)

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
	r.debug.Draw(r.camera)
	r.camera.DrawBuffer(screen)
}

func (r *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.ScreenWidth * common.Scale, common.ScreenHeight * common.Scale
}
