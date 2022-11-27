package core

import (
	"platformer/common"
)

type Level struct {
	name      string
	tiledGrid *common.TiledGrid
}

func NewLevel(name string) *Level {
	l := &Level{}
	l.tiledGrid = common.NewTileGrid(name)
	return l
}

func (r *Level) Update(delta float64, game *Game) {

}

func (r *Level) Draw(camera common.Camera) {
	r.tiledGrid.Draw(camera)
}
