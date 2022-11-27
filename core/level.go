package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"platformer/common"
)

type Level struct {
	name             string
	tiledGrid        *common.TiledGrid
	background       *ebiten.Image
	backgroundOffset float64
}

func NewLevel(name string) *Level {
	l := &Level{
		background:       common.LoadImage(name + "/background.png"),
		backgroundOffset: 60,
	}
	l.tiledGrid = common.NewTileGrid(name)
	return l
}

func (r *Level) Update(delta float64, game *Game) {
	r.backgroundOffset = (game.player.y / float64(r.tiledGrid.Layers[0].Height*common.TileSize)) * (60)
}

func (r *Level) Draw(camera common.Camera) {

	cx, cy := camera.GetPos()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cx, cy-r.backgroundOffset)
	op.GeoM.Scale(common.Scale, common.Scale)
	camera.DrawImage(r.background, op)

	r.tiledGrid.Draw(camera)
}
