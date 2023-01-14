package core

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"platformer/common"
)

type Pickup struct {
	x      float64
	y      float64
	image  *ebiten.Image
	effect Effect
	done   bool
}

func (r *Pickup) Update(delta float64, game *Game) {
	if common.Overlap(game.Player.x+4, game.Player.y+8, 8, 8, r.x+2, r.y+2, 12, 12) {
		r.effect.GetPickedUp(game)
		game.Level.RemovePickup(r)
	}
}

func (r *Pickup) Draw(camera common.Camera) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	camera.DrawImage(r.image, op)
}

type Effect interface {
	GetPickedUp(game *Game)
}

type HealthEffect struct {
	amount int
}

func (r *HealthEffect) GetPickedUp(game *Game) {
	if game.Player.Health < game.Player.MaxHealth {
		game.Player.AddHealth(1)
	}
}

type BookEffect struct {
	title string
	spell string
}

func (r *BookEffect) GetPickedUp(game *Game) {
	fmt.Println("picked up the book called: ", r.title)
	game.Player.AddSpell(r.spell)
	game.PlayerProgress.AddSpell(r.spell)
}
