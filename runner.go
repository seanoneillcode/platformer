package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"platformer/common"
	"platformer/core"
	"platformer/gui"
	"platformer/res"
	"time"
)

type Runner struct {
	lastUpdateCalled time.Time
	firstUpdate      bool

	// refs
	res           *res.Resources
	game          *core.Game
	userInterface *gui.UserInterface
}

func NewRunner() *Runner {
	resources := res.NewResources()
	r := &Runner{
		firstUpdate:   true,
		res:           resources,
		userInterface: gui.NewUserInterface(resources),
	}
	r.game = core.NewGame(resources, r)
	return r
}

func (r *Runner) Update() error {
	if r.firstUpdate {
		// we want a reasonable value for the delta, so we 'skip' the first update,
		// otherwise the first movement of everything that moves is scaled by a large delta
		// causing things to go through walls.
		r.firstUpdate = false
		r.lastUpdateCalled = time.Now()
		return nil
	}
	delta := float64(time.Now().Sub(r.lastUpdateCalled).Milliseconds()) / 1000
	r.lastUpdateCalled = time.Now()

	err := r.game.Update(delta)
	if err != nil {
		return err
	}

	err = r.userInterface.Update(delta, r.game)
	if err != nil {
		return err
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return common.NormalEscapeError
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	return nil
}

func (r *Runner) Draw(screen *ebiten.Image) {
	r.game.Draw(screen)
	r.userInterface.Draw(screen)
}

func (r *Runner) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.ScreenWidth * common.Scale, common.ScreenHeight * common.Scale
}

func (r *Runner) CloseBook() {
	r.game.Enabled = true
	r.userInterface.Enabled = false
}

func (r *Runner) OpenBook(title, text string) {
	r.game.Enabled = false
	r.userInterface.Enabled = true
	r.userInterface.Book.Open(title, text)
}
