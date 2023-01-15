package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"platformer/common"
)

const fadeInTimeAmount = 0.6
const fadeOutTimeAmount = 0.8

type Sign struct {
	x          float64
	y          float64
	image      *ebiten.Image
	popupImage *ebiten.Image
	text       string

	// popup
	px float64
	py float64
	pw float64

	// state
	wasVisible bool
	isVisible  bool
	fadeAmount float64
	timer      float64
}

func NewSign(x, y float64, text string, game *Game) *Sign {
	pw := common.TextWidth(text) + 4 + 4
	px := x - (float64(pw) / 2.0) + 8
	py := y - 24
	return &Sign{
		x:          x,
		y:          y,
		text:       text,
		image:      game.res.GetImage("sign"),
		popupImage: game.res.GetImage("popup-sign"),
		pw:         float64(pw),
		px:         px,
		py:         py,
		timer:      -1,
	}
}

func (r *Sign) Update(delta float64, game *Game) {

	overlap := common.Overlap(game.Player.x, game.Player.y, game.Player.sizex, game.Player.sizey, r.x, r.y, 16, 16)

	if overlap {
		r.isVisible = true
		if !r.wasVisible {
			if r.timer < 0 {
				r.timer = fadeInTimeAmount
			}
			r.wasVisible = true
		}
	} else {
		if r.wasVisible {
			r.wasVisible = false
			if r.timer < 0 {
				r.timer = fadeOutTimeAmount
			}
		}
	}

	if r.timer > 0 {
		r.timer = r.timer - delta
		if overlap {
			r.fadeAmount = 1.0 - r.timer
		} else {
			r.fadeAmount = r.timer
		}
		if r.timer < 0 {
			if !r.wasVisible {
				r.isVisible = false
			}
		}
	}

	r.py = r.y - 24 - (r.fadeAmount * 8)

}

func (r *Sign) Draw(camera common.Camera) {
	// draw sign in the level
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	camera.DrawImage(r.image, op)

	// draw text pop up
	if !r.isVisible {
		return
	}
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.px, r.py)
	op.GeoM.Scale(common.Scale, common.Scale)
	op.ColorM.Scale(1, 1, 1, r.fadeAmount)
	img := r.popupImage.SubImage(image.Rect(0, 0, int(r.pw), 14)).(*ebiten.Image)
	camera.DrawImage(img, op)

	common.DrawTextWithAlpha(camera, r.text, r.px+4, r.py+4, r.fadeAmount)
}
