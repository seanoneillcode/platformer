package common

import "github.com/hajimehoshi/ebiten/v2"

type Camera interface {
	DrawImage(img *ebiten.Image, options *ebiten.DrawImageOptions)
	GetPos() (float64, float64)
}
