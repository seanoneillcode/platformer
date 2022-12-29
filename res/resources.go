package res

import (
	"github.com/hajimehoshi/ebiten/v2"
	"platformer/common"
)

type Resources struct {
	images map[string]*ebiten.Image
}

func NewResources() *Resources {
	return &Resources{
		images: map[string]*ebiten.Image{
			"player-run":            common.LoadImage("player-run.png"),
			"player-idle":           common.LoadImage("player-idle.png"),
			"player-jump":           common.LoadImage("player-jump.png"),
			"player-fall":           common.LoadImage("player-fall.png"),
			"player-hurt":           common.LoadImage("player-hurt.png"),
			"player-death":          common.LoadImage("player-death.png"),
			"player-climb":          common.LoadImage("player-climb.png"),
			"player-crouch":         common.LoadImage("player-crouch.png"),
			"book-pickup":           common.LoadImage("book.png"),
			"health-pickup":         common.LoadImage("health.png"),
			"crawler-run":           common.LoadImage("crawler-run.png"),
			"crawler-idle":          common.LoadImage("crawler-idle.png"),
			"crawler-hurt":          common.LoadImage("crawler-hurt.png"),
			"crawler-die":           common.LoadImage("crawler-die.png"),
			"blob-run":              common.LoadImage("blob-run.png"),
			"blob-idle":             common.LoadImage("blob-idle.png"),
			"blob-hurt":             common.LoadImage("blob-hurt.png"),
			"blob-die":              common.LoadImage("blob-die.png"),
			"blob-attack":           common.LoadImage("blob-attack.png"),
			"spell-bullet":          common.LoadImage("spell-bullet.png"),
			"effect-spell-hit":      common.LoadImage("effect-spell-hit.png"),
			"health-bar":            common.LoadImage("health-bar.png"),
			"health-bar-background": common.LoadImage("health-bar-background.png"),
			"health-bar-end":        common.LoadImage("health-bar-end.png"),
			"flimsy":                common.LoadImage("flimsy.png"),
		},
	}
}

func (r *Resources) GetImage(name string) *ebiten.Image {
	img, ok := r.images[name]
	if !ok {
		panic("missing resource " + name)
	}
	return img
}
