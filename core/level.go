package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"platformer/common"
)

const spawnObject = "spawn"
const exitObject = "exit"

type Level struct {
	name             string
	tiledGrid        *common.TiledGrid
	background       *ebiten.Image
	backgroundOffset float64
	spawn            *Spawn
	exit             *Exit
}

func NewLevel(name string) *Level {
	l := &Level{
		background:       common.LoadImage(name + "/background.png"),
		backgroundOffset: 60,
	}
	l.tiledGrid = common.NewTileGrid(name)
	objects := l.tiledGrid.GetObjectData()
	for _, object := range objects {
		if object.Name == spawnObject {
			l.spawn = &Spawn{
				x: float64(object.X),
				y: float64(object.Y),
			}
		}
		if object.Name == exitObject {
			l.exit = &Exit{
				x: float64(object.X),
				y: float64(object.Y),
			}
			for _, prop := range object.Properties {
				if prop.Name == "next-level" && prop.Value != nil {
					l.exit.nextLevel = (prop.Value).(string)
				}
			}
		}
	}
	// validate level
	if l.spawn == nil {
		panic("no spawn for level: " + name)
	}
	return l
}

func (r *Level) Update(delta float64, game *Game) {
	r.backgroundOffset = (game.player.y / float64(r.tiledGrid.Layers[0].Height*common.TileSize)) * (60)
	if r.exit != nil {
		if common.Overlap(game.player.x+8, game.player.y+4, game.player.sizex, game.player.sizey, r.exit.x, r.exit.y, common.TileSize, common.TileSize*2) {
			game.MoveToNextLevel(r.exit.nextLevel)
		}
	}
}

func (r *Level) Draw(camera common.Camera) {

	cx, cy := camera.GetPos()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cx, cy-r.backgroundOffset)
	op.GeoM.Scale(common.Scale, common.Scale)
	camera.DrawImage(r.background, op)

	r.tiledGrid.Draw(camera)
}

type Spawn struct {
	x float64
	y float64
}

type Exit struct {
	x         float64
	y         float64
	nextLevel string
}
