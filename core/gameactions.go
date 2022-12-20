package core

import "fmt"

func (r *Game) LoadLevel(name string) {
	r.player = NewPlayer(r)
	r.camera = NewCamera()
	r.camera.Target(r.player)
	r.level = NewLevel(name, r)
	r.player.x = r.level.spawn.x
	r.player.y = r.level.spawn.y
	fmt.Println("load level ", name)
}

func (r *Game) PlayerDeath() {
	r.player = NewPlayer(r)
	r.camera.Target(r.player)
	r.player.x = r.level.spawn.x
	r.player.y = r.level.spawn.y
	fmt.Println("player death")
}

func (r *Game) MoveToNextLevel(level string) {
	r.LoadLevel(level)
}
