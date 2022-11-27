package core

func (r *Game) LoadLevel(name string) {
	r.player = NewPlayer(r)
	r.camera = NewCamera()
	r.camera.Target(r.player)
	r.level = NewLevel(name)
}
