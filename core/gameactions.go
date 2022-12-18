package core

import "fmt"

func (r *Game) LoadLevel(name string) {
	r.player = NewPlayer(r)
	r.camera = NewCamera()
	r.camera.Target(r.player)
	r.level = NewLevel(name)
	r.player.x = r.level.spawn.x
	r.player.y = r.level.spawn.y
	fmt.Println("spawn ", r.player.x, " ", r.player.y)
}
