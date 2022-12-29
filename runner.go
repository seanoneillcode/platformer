package main

import (
	"github.com/hajimehoshi/ebiten/v2"
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
	return &Runner{
		firstUpdate:   true,
		res:           resources,
		game:          core.NewGame(resources),
		userInterface: gui.NewUserInterface(resources),
	}
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

	return nil
}

func (r *Runner) Draw(screen *ebiten.Image) {
	r.game.Draw(screen)
	r.userInterface.Draw(screen)
}

func (r *Runner) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.ScreenWidth * common.Scale, common.ScreenHeight * common.Scale
}
