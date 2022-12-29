package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"platformer/core"
	"platformer/res"
)

type UserInterface struct {
	hud *Hud
}

func NewUserInterface(resources *res.Resources) *UserInterface {
	return &UserInterface{
		hud: NewHud(resources),
	}
}

func (r *UserInterface) Update(delta float64, game *core.Game) error {
	r.hud.Update(delta, game)
	return nil
}

func (r *UserInterface) Draw(screen *ebiten.Image) {
	r.hud.Draw(screen)
}
