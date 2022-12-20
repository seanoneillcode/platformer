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
	pickups          []*Pickup
}

func NewLevel(name string, game *Game) *Level {
	l := &Level{
		background:       common.LoadImage(name + "/background.png"),
		backgroundOffset: 60,
	}
	l.tiledGrid = common.NewTileGrid(name)
	objects := l.tiledGrid.GetObjectData()
	l.pickups = []*Pickup{}
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
		if object.Name == healthPickup {
			newPickup := &Pickup{
				x:     float64(object.X),
				y:     float64(object.Y),
				image: game.images["health-pickup"],
			}
			effect := &HealthEffect{
				amount: 1,
			}
			for _, prop := range object.Properties {
				if prop.Name == "amount" && prop.Value != nil {
					effect.amount = (prop.Value).(int)
				}
			}
			newPickup.effect = effect
			l.pickups = append(l.pickups, newPickup)
		}
		if object.Name == bookPickup {
			newPickup := &Pickup{
				x:     float64(object.X),
				y:     float64(object.Y),
				image: game.images["book-pickup"],
			}
			effect := &BookEffect{
				title: "untitled",
			}
			for _, prop := range object.Properties {
				if prop.Name == "title" && prop.Value != nil {
					effect.title = (prop.Value).(string)
				}
			}
			newPickup.effect = effect
			l.pickups = append(l.pickups, newPickup)
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
	for _, pickup := range r.pickups {
		pickup.Update(delta, game)
	}
}

func (r *Level) Draw(camera common.Camera) {

	cx, cy := camera.GetPos()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cx, cy-r.backgroundOffset)
	op.GeoM.Scale(common.Scale, common.Scale)
	camera.DrawImage(r.background, op)

	r.tiledGrid.Draw(camera)

	for _, pickup := range r.pickups {
		pickup.Draw(camera)
	}
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

func (r *Level) RemovePickup(pickup *Pickup) {
	newListOfPickups := []*Pickup{}
	for _, p := range r.pickups {
		if p != pickup {
			newListOfPickups = append(newListOfPickups, p)
		}
	}
	r.pickups = newListOfPickups
}
