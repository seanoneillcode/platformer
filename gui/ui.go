package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"platformer/core"
	"platformer/res"
)

type UserInterface struct {
	hud     *Hud
	Book    *Book
	Enabled bool
}

func NewUserInterface(resources *res.Resources) *UserInterface {
	return &UserInterface{
		hud:  NewHud(resources),
		Book: NewBook(resources),
	}
}

func (r *UserInterface) Update(delta float64, game *core.Game) error {
	r.hud.Update(delta, game)
	r.Book.Update(delta, game)
	return nil
}

func (r *UserInterface) Draw(screen *ebiten.Image) {
	r.hud.Draw(screen)
	r.Book.Draw(screen)
}
