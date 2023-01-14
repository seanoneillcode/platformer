package core

import "fmt"

func (r *Game) LoadLevel(name string) {
	r.spellObjects = []*SpellObject{}
	r.effectSprites = []*EffectSprite{}
	r.Level = NewLevel(name, r)
	r.Player = NewPlayer(r)
	r.Player.x = r.Level.spawn.x
	r.Player.y = r.Level.spawn.y
	r.PlayerProgress.HydratePlayer(r.Player)
	r.Camera = NewCamera()
	r.Camera.Target(r.Player)
	fmt.Println("load Level ", name)
}

func (r *Game) PlayerDeath() {
	r.spellObjects = []*SpellObject{}
	r.effectSprites = []*EffectSprite{}
	r.Player = NewPlayer(r)
	r.Player.x = r.Level.spawn.x
	r.Player.y = r.Level.spawn.y
	r.PlayerProgress.HydratePlayer(r.Player)
	r.Camera.Target(r.Player)
	fmt.Println("Player death")
}

func (r *Game) MoveToNextLevel(level string) {
	r.LoadLevel(level)
}
