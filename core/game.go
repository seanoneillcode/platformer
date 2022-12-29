package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"platformer/common"
	"platformer/res"
)

type Game struct {
	Player        *Player
	Camera        *Camera
	Level         *Level
	debug         *DebugDrawer
	spellObjects  []*SpellObject
	effectSprites []*EffectSprite
	// refs
	res *res.Resources
}

func NewGame(resources *res.Resources) *Game {
	r := &Game{
		debug: NewDebug(),
		res:   resources,
	}
	r.LoadLevel("level-alpha")
	return r
}

func (r *Game) Update(delta float64) error {
	r.debug.Update(delta, r)
	r.Player.Update(delta, r)
	r.Level.Update(delta, r)
	r.Camera.Update(delta, r)
	for _, s := range r.spellObjects {
		s.Update(delta, r)
	}
	for _, e := range r.effectSprites {
		e.Update(delta, r)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return common.NormalEscapeError
	}

	return nil
}

func (r *Game) Draw(screen *ebiten.Image) {
	r.Level.Draw(r.Camera)
	r.Player.Draw(r.Camera)
	for _, s := range r.spellObjects {
		s.Draw(r.Camera)
	}
	for _, e := range r.effectSprites {
		e.Draw(r.Camera)
	}
	r.debug.Draw(r.Camera)
	r.Camera.DrawBuffer(screen)
}

func (r *Game) RemoveSpellObject(spellObject *SpellObject) {
	newSpellObjs := []*SpellObject{}
	for _, so := range r.spellObjects {
		if so != spellObject {
			newSpellObjs = append(newSpellObjs, so)
		}
	}
	r.spellObjects = newSpellObjs
}

func (r *Game) AddSpellObject(spellObject *SpellObject) {
	r.spellObjects = append(r.spellObjects, spellObject)
}

func (r *Game) RemoveEffectSprite(effectSprite *EffectSprite) {
	newEffectSprites := []*EffectSprite{}
	for _, so := range r.effectSprites {
		if so != effectSprite {
			newEffectSprites = append(newEffectSprites, so)
		}
	}
	r.effectSprites = newEffectSprites
}

func (r *Game) AddEffectSprite(effectSprite *EffectSprite) {
	r.effectSprites = append(r.effectSprites, effectSprite)
}

const (
	effectSpellHit     = "effect-spell-hit"
	effectCrawlerDeath = "effect-crawler-death"
	effectBlobDeath    = "effect-blob-death"
)

func (r *Game) SpawnEffect(name string, x, y float64, isFlip bool) {
	switch name {
	case effectCrawlerDeath:
		r.AddEffectSprite(&EffectSprite{
			x: x,
			y: y,
			w: 24,
			h: 24,
			animation: &Animation{
				image:           r.res.GetImage("crawler-die"),
				numFrames:       2,
				size:            24,
				frameTimeAmount: 0.2,
				isLoop:          false,
			},
			isFlip: isFlip,
		})
	case effectBlobDeath:
		r.AddEffectSprite(&EffectSprite{
			x: x,
			y: y,
			w: 32,
			h: 32,
			animation: &Animation{
				image:           r.res.GetImage("blob-die"),
				numFrames:       6,
				size:            32,
				frameTimeAmount: 0.1,
				isLoop:          false,
			},
			isFlip: isFlip,
		})
	case effectSpellHit:
		r.AddEffectSprite(&EffectSprite{
			x: x - 4,
			y: y - 4,
			w: 24,
			h: 24,
			animation: &Animation{
				image:           r.res.GetImage("effect-spell-hit"),
				numFrames:       5,
				size:            24,
				frameTimeAmount: 0.1,
				isLoop:          false,
			},
			isTemporary: true,
			ttl:         0.4,
			isFlip:      isFlip,
		})

	}
}
