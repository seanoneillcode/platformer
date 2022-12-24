package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"platformer/common"
)

type Hud struct {
	healthBarBackgroundImage *ebiten.Image
	healthBarHealthImage     *ebiten.Image
	healthBarEndImage        *ebiten.Image
	healthPercent            float64
}

func NewHud(game *Game) *Hud {
	return &Hud{
		healthBarBackgroundImage: game.images["health-bar-background"],
		healthBarEndImage:        game.images["health-bar-end"],
		healthBarHealthImage:     game.images["health-bar"],
	}
}

func (r *Hud) Update(delta float64, game *Game) {
	// insert logic to track health here
	r.healthPercent = float64(game.player.health) / float64(game.player.maxHealth)
}

func (r *Hud) Draw(camera common.Camera) {

	healthPercentXPos := 46 * r.healthPercent

	cx, cy := camera.GetPos()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cx+4, cy+4)
	op.GeoM.Scale(common.Scale, common.Scale)
	camera.DrawImage(r.healthBarBackgroundImage, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cx+4, cy+4)
	op.GeoM.Scale(common.Scale, common.Scale)
	img := r.healthBarHealthImage.SubImage(image.Rect(0, 0, int(healthPercentXPos), 6)).(*ebiten.Image)
	camera.DrawImage(img, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cx+3+healthPercentXPos, cy+5)
	op.GeoM.Scale(common.Scale, common.Scale)
	camera.DrawImage(r.healthBarEndImage, op)
}
