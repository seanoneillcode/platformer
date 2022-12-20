package core

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"platformer/common"
)

const (
	healthPickup = "health"
	bookPickup   = "book"
)

type Pickup struct {
	x      float64
	y      float64
	image  *ebiten.Image
	effect Effect
	done   bool
}

func (r *Pickup) Update(delta float64, game *Game) {
	if common.Overlap(game.player.x+4, game.player.y+8, 8, 8, r.x+2, r.y+2, 12, 12) {
		r.effect.GetPickedUp(game)
		game.level.RemovePickup(r)
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
	if game.player.health < game.player.maxHealth {
		game.player.AddHealth(1)
	}
}

type BookEffect struct {
	title string
}

func (r *BookEffect) GetPickedUp(game *Game) {
	fmt.Println("picked up the book called: ", r.title)
}
