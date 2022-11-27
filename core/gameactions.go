package core

func (r *Game) StartNewGame() {
	r.player = NewPlayer(r)
}
