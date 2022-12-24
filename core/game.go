package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"platformer/common"
	"time"
)

type Game struct {
	lastUpdateCalled time.Time
	player           *Player
	images           map[string]*ebiten.Image
	camera           *Camera
	level            *Level
	debug            *DebugDrawer
	spellObjects     []*SpellObject
	effectSprites    []*EffectSprite
}

func NewGame() *Game {
	r := &Game{
		images: map[string]*ebiten.Image{
			"player-run":       common.LoadImage("player-run.png"),
			"player-idle":      common.LoadImage("player-idle.png"),
			"player-jump":      common.LoadImage("player-jump.png"),
			"player-fall":      common.LoadImage("player-fall.png"),
			"player-hurt":      common.LoadImage("player-hurt.png"),
			"player-death":     common.LoadImage("player-death.png"),
			"player-climb":     common.LoadImage("player-climb.png"),
			"player-crouch":    common.LoadImage("player-crouch.png"),
			"book-pickup":      common.LoadImage("book.png"),
			"health-pickup":    common.LoadImage("health.png"),
			"crawler-run":      common.LoadImage("crawler-run.png"),
			"crawler-idle":     common.LoadImage("crawler-idle.png"),
			"crawler-hurt":     common.LoadImage("crawler-hurt.png"),
			"crawler-die":      common.LoadImage("crawler-die.png"),
			"spell-bullet":     common.LoadImage("spell-bullet.png"),
			"effect-spell-hit": common.LoadImage("effect-spell-hit.png"),
		},
		lastUpdateCalled: time.Now(),
		debug:            NewDebug(),
	}
	r.LoadLevel("enemy-level")
	return r
}

func (r *Game) Update() error {
	delta := float64(time.Now().Sub(r.lastUpdateCalled).Milliseconds()) / 1000
	r.lastUpdateCalled = time.Now()

	r.debug.Update(delta, r)
	r.player.Update(delta, r)
	r.level.Update(delta, r)
	r.camera.Update(delta, r)
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
	r.level.Draw(r.camera)
	r.player.Draw(r.camera)
	for _, s := range r.spellObjects {
		s.Draw(r.camera)
	}
	for _, e := range r.effectSprites {
		e.Draw(r.camera)
	}
	r.debug.Draw(r.camera)
	r.camera.DrawBuffer(screen)
}

func (r *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.ScreenWidth * common.Scale, common.ScreenHeight * common.Scale
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
				image:           r.images["crawler-die"],
				numFrames:       2,
				size:            24,
				frameTimeAmount: 0.2,
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
				image:           r.images["effect-spell-hit"],
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
