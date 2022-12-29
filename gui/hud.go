package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"platformer/common"
	"platformer/core"
	"platformer/res"
)

type Hud struct {
	healthBarBackgroundImage *ebiten.Image
	healthBarHealthImage     *ebiten.Image
	healthBarEndImage        *ebiten.Image
	healthPercent            float64
}

func NewHud(resources *res.Resources) *Hud {
	return &Hud{
		healthBarBackgroundImage: resources.GetImage("health-bar-background"),
		healthBarEndImage:        resources.GetImage("health-bar-end"),
		healthBarHealthImage:     resources.GetImage("health-bar"),
	}
}

func (r *Hud) Update(delta float64, game *core.Game) {
	r.healthPercent = float64(game.Player.Health) / float64(game.Player.MaxHealth)
}

func (r *Hud) Draw(screen *ebiten.Image) {

	healthPercentXPos := 46 * r.healthPercent

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(4, 4)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(r.healthBarBackgroundImage, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(4, 4)
	op.GeoM.Scale(common.Scale, common.Scale)
	img := r.healthBarHealthImage.SubImage(image.Rect(0, 0, int(healthPercentXPos), 6)).(*ebiten.Image)
	screen.DrawImage(img, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(3+healthPercentXPos, 5)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(r.healthBarEndImage, op)
}
