package core

type PlayerProgress struct {
	mostRecentSpell string
	spells          map[string]bool
}

func (r *PlayerProgress) AddSpell(spell string) {
	_, ok := r.spells[spell]
	if !ok {
		r.spells[spell] = true
	}
	if len(r.spells) == 1 {
		r.mostRecentSpell = spell
	}
}

func (r *PlayerProgress) HydratePlayer(player *Player) {
	player.currentSpell = r.mostRecentSpell
	player.spells = r.spells
}
