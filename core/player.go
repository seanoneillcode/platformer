package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"math"
	"platformer/common"
)

const playingState = "playing"
const dyingState = "dying"

const standardJumpHeight = 16 * 3.2
const forcedJumpHeight = 16 * 2
const standardJumpTime = 0.44
const standardFallTime = 0.4
const minimumJumpHeight = 16
const coyoteTimeAmount = 0.16
const fudge = 0.001
const runAcc = 20.0
const maxRunVelocity = 100
const ladderVelocity = 70
const lateJumpMarginTime = 0.14
const ladderGrabAllowance = 8.0
const takeDamageTime = 0.3
const postDamageTime = 0.6
const playerDeathTime = 2.0
const castSpellCoolDownTime = 0.2
const castSpellTimeTotal = 0.3

type Player struct {
	x                  float64
	y                  float64
	targetY            float64
	sizex              float64
	sizey              float64
	drawOffsetX        float64
	drawOffsetY        float64
	drawSizex          int
	drawSizey          int
	state              string
	timer              float64
	shootTimer         float64
	image              *ebiten.Image
	animTimer          float64
	targetFrame        int
	moveYSpeed         float64
	jumpAcc            float64
	velocityY          float64
	velocityX          float64
	coyoteTimer        float64
	jumpTimer          float64
	wasPressingJump    bool
	alreadyAbortedJump bool
	targetVelocityX    float64
	isFlip             bool
	currentAnimation   string
	animations         map[string]*Animation
	lateJumpTimer      float64
	lockedToLadder     bool
	takeDamageTimer    float64
	postDamageTimer    float64
	Health             int
	deathTimer         float64
	MaxHealth          int
	castSpellTimer     float64
	isCrouch           bool
	currentSpell       string
	spells             map[string]bool
	castTimer          float64
}

func NewPlayer(game *Game) *Player {
	p := &Player{
		state:            playingState,
		Health:           6,
		MaxHealth:        9,
		x:                19 * common.TileSize,
		y:                12 * common.TileSize,
		sizex:            16, // physical size
		sizey:            16, // physical size
		drawOffsetX:      8,  // just for drawing
		drawOffsetY:      8,  // just for drawing
		drawSizex:        32, // just for drawing
		drawSizey:        32,
		currentAnimation: "idle",
		spells:           map[string]bool{},
		currentSpell:     "spell-bullet", // spell-bullet
		animations: map[string]*Animation{
			"run": {
				image:           game.res.GetImage("player-run"),
				numFrames:       6,
				size:            32,
				frameTimeAmount: 0.1,
				isLoop:          true,
			},
			"run-cast": {
				image:           game.res.GetImage("player-run-cast"),
				numFrames:       6,
				size:            32,
				frameTimeAmount: 0.1,
				isLoop:          true,
			},
			"idle": {
				image:           game.res.GetImage("player-idle"),
				numFrames:       1,
				size:            32,
				frameTimeAmount: 1,
				isLoop:          true,
			},
			"crouch": {
				image:           game.res.GetImage("player-crouch"),
				numFrames:       1,
				size:            32,
				frameTimeAmount: 1,
				isLoop:          true,
			},
			"jump": {
				image:           game.res.GetImage("player-jump"),
				numFrames:       1,
				size:            32,
				frameTimeAmount: 1,
				isLoop:          true,
			},
			"jump-cast": {
				image:           game.res.GetImage("player-jump-cast"),
				numFrames:       1,
				size:            32,
				frameTimeAmount: 1,
				isLoop:          true,
			},
			"death": {
				image:           game.res.GetImage("player-death"),
				numFrames:       3,
				size:            32,
				frameTimeAmount: 0.4,
				isLoop:          false,
			},
			"fall": {
				image:           game.res.GetImage("player-fall"),
				numFrames:       1,
				size:            32,
				frameTimeAmount: 1,
				isLoop:          true,
			},
			"fall-cast": {
				image:           game.res.GetImage("player-fall-cast"),
				numFrames:       1,
				size:            32,
				frameTimeAmount: 1,
				isLoop:          true,
			},
			"climb": {
				image:           game.res.GetImage("player-climb"),
				numFrames:       2,
				size:            32,
				frameTimeAmount: 0.2,
				isLoop:          true,
			},
			"hurt": {
				image:           game.res.GetImage("player-hurt"),
				numFrames:       1,
				size:            32,
				frameTimeAmount: 1,
				isLoop:          true,
			},
			"cast": {
				image:           game.res.GetImage("player-cast"),
				numFrames:       3,
				size:            32,
				frameTimeAmount: 0.1,
				isLoop:          false,
			},
		},
		velocityY:       0,
		velocityX:       0,
		targetVelocityX: 0,
	}
	return p
}

func (r *Player) Update(delta float64, game *Game) {
	var aimY float64

	switch r.state {
	case dyingState:
		r.deathTimer -= delta
		if r.deathTimer < 0 {
			game.PlayerDeath()
		}
		r.currentAnimation = "death"
		r.animations[r.currentAnimation].Update(delta)
	case playingState:
		var tryJump bool
		var pressJump bool
		var tryFall bool
		var tryMoveY = 0.0
		r.targetVelocityX = 0
		r.currentAnimation = "idle"
		shouldUpdateAnimation := false

		if r.takeDamageTimer <= 0 {
			if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
				r.targetVelocityX = -maxRunVelocity
				r.isFlip = true
				r.currentAnimation = "run"
				if r.lockedToLadder {
					r.targetVelocityX = -maxRunVelocity / 2.0
				}
				shouldUpdateAnimation = true
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyC) {
				game.Actions.OpenBook("my title", "My friend,\n\n\nI have fallen and can't get up.\n\nCan you help me?\n\nI can write a PROPER sentence now, full of lore and \nsuspense.")
			}
			if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
				r.targetVelocityX = maxRunVelocity
				r.isFlip = false
				r.currentAnimation = "run"
				if r.lockedToLadder {
					r.targetVelocityX = maxRunVelocity / 2.0
				}
				shouldUpdateAnimation = true
			}
			r.isCrouch = false
			if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
				tryFall = true
				tryMoveY = 1
				shouldUpdateAnimation = true
				if r.targetVelocityX == 0 {
					r.currentAnimation = "crouch"
				}
				r.isCrouch = true
			}
			if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
				tryMoveY = -1
				shouldUpdateAnimation = true
			}

			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				pressJump = true
				shouldUpdateAnimation = true
			}

			if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
				tryJump = true
				pressJump = true
				r.lateJumpTimer = lateJumpMarginTime
				shouldUpdateAnimation = true
			}
		} else {
			r.targetVelocityX = maxRunVelocity / 2.0
			if r.isFlip {
				r.targetVelocityX = r.targetVelocityX * -1
			}
		}

		r.lateJumpTimer = r.lateJumpTimer - delta
		time := standardJumpTime
		if r.jumpTimer > (standardJumpTime) {
			time = standardFallTime
		}
		gravity := (standardJumpHeight * -2) / (time * time)

		oldx := r.x
		oldy := r.y

		partial := 4.0

		newx := r.x + (delta * r.velocityX)
		movey := 0.0
		if !r.lockedToLadder {
			movey = (r.velocityY * delta) + (0.5 * gravity * delta * delta)
		}
		newy := r.y - movey
		r.velocityY = r.velocityY + (gravity * delta)

		cr := DoCollision(oldx+partial, oldy+partial, newx+partial, newy+partial, r.sizex-partial-partial, r.sizey+partial, game.Level, tryFall)

		newx = cr.newX - partial
		newy = cr.newY - partial

		if cr.hitFloor {
			r.velocityY = 0
			r.coyoteTimer = coyoteTimeAmount
		}

		var touchingLadder = false
		tx, ty := int((oldx+(r.sizex/2.0))/common.TileSize), int((oldy+partial)/common.TileSize)
		td := game.Level.tiledGrid.GetTileData(tx, ty)
		if td.Ladder {
			touchingLadder = true
			if tryMoveY != 0 {
				middle := float64(tx * common.TileSize)
				left, right := middle-ladderGrabAllowance, middle+ladderGrabAllowance
				if oldx > left && oldx < right {
					r.lockedToLadder = true
					newy = oldy + (delta * ladderVelocity * tryMoveY)
				}
			}

		}
		tx, ty = int((oldx+(r.sizex/2.0))/common.TileSize), int((oldy+r.sizey+partial+partial)/common.TileSize)
		td = game.Level.tiledGrid.GetTileData(tx, ty)
		if td.Ladder {
			touchingLadder = true
			if tryMoveY != 0 {
				middle := float64(tx * common.TileSize)
				left, right := middle-ladderGrabAllowance, middle+ladderGrabAllowance
				if oldx > left && oldx < right {
					r.lockedToLadder = true
					newy = oldy + (delta * ladderVelocity * tryMoveY)

					// if move down, check for block
					if tryMoveY == 1 {
						tx, ty = int((oldx+(r.sizex/2.0))/common.TileSize), int((newy+r.sizey+partial+partial)/common.TileSize)
						td = game.Level.tiledGrid.GetTileData(tx, ty)
						if td.Block {
							r.lockedToLadder = false
							newy = oldy
						}
					}
				}
			}
		}
		if !touchingLadder {
			r.lockedToLadder = false
		}
		var hitDamage bool
		tx, ty = int((oldx+(r.sizex/2.0))/common.TileSize), int((oldy+r.sizey)/common.TileSize)
		td = game.Level.tiledGrid.GetTileData(tx, ty)
		if td.Damage {
			hitDamage = true
		}

		r.x = newx
		r.y = newy

		if tryJump || r.lateJumpTimer > 0 {
			if cr.hitFloor || r.coyoteTimer > 0 || r.lockedToLadder {
				r.lockedToLadder = false
				r.coyoteTimer = 0
				r.velocityY = (2 * standardJumpHeight) / standardJumpTime
				r.jumpTimer = 0
				r.alreadyAbortedJump = false
			}
		}
		if !pressJump {
			// if Player is currently jumping in the first half phase of jumping
			if r.jumpTimer < (standardJumpTime*0.5) && r.wasPressingJump && !r.alreadyAbortedJump {
				r.alreadyAbortedJump = true
				r.velocityY = (2 * minimumJumpHeight) / (standardJumpTime)
			}
		}
		r.jumpTimer = r.jumpTimer + delta
		if r.coyoteTimer > 0 {
			r.coyoteTimer = r.coyoteTimer - delta
		}
		if cr.hitFloor {
			r.lockedToLadder = false
		}

		if oldy != newy {
			if r.velocityY <= 0 {
				r.currentAnimation = "fall"
			}
			if r.velocityY > 0 {
				r.currentAnimation = "jump"
			}
		}
		if r.lockedToLadder {
			r.currentAnimation = "climb"
			r.velocityY = 0
			r.alreadyAbortedJump = false
		}
		if cr.hitCeiling {
			r.alreadyAbortedJump = true
			r.velocityY = -20
		}
		r.wasPressingJump = pressJump
		if cr.hitWall {
			r.targetVelocityX = 0
			r.velocityX = 0
		}
		if (tryMoveY == 1 && !cr.hitFloor) || tryMoveY == -1 {
			aimY = tryMoveY
		}

		if r.velocityX < r.targetVelocityX {
			r.velocityX = r.velocityX + runAcc
			if r.velocityX > r.targetVelocityX {
				r.velocityX = r.targetVelocityX
			}
		}
		if r.velocityX > r.targetVelocityX {
			r.velocityX = r.velocityX - runAcc
			if r.velocityX < r.targetVelocityX {
				r.velocityX = r.targetVelocityX
			}
		}
		if r.castTimer > 0 {
			r.castTimer = r.castTimer - delta
			switch r.currentAnimation {
			case "run":
				r.currentAnimation = "run-cast"
				r.animations["run"].Update(delta)
				shouldUpdateAnimation = true
			case "jump":
				r.currentAnimation = "jump-cast"
			case "fall":
				r.currentAnimation = "fall-cast"
			default:
				shouldUpdateAnimation = true
				r.currentAnimation = "cast"
			}

		}
		if hitDamage {
			game.Player.TakeDamage(game)
		}
		if r.takeDamageTimer > 0 {
			r.takeDamageTimer -= delta
			if r.takeDamageTimer < 0 {
				r.postDamageTimer = postDamageTime
			}
			r.currentAnimation = "hurt"
		}
		if r.postDamageTimer > 0 {
			r.postDamageTimer -= delta
			if r.postDamageTimer < 0 {
				// normal
			}
		}

		if shouldUpdateAnimation {
			r.animations[r.currentAnimation].Update(delta)
		}
	}

	r.castSpellTimer = r.castSpellTimer - delta
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if r.castSpellTimer < 0 {
			switch r.currentSpell {
			case "spell-bullet":
				r.castSpellTimer = castSpellCoolDownTime
				r.castTimer = castSpellTimeTotal
				r.animations["cast"].Reset()
				r.animations["run-cast"].SnapToAnimation(r.animations["run"])

				var moveX float64
				var moveY float64
				var posX float64
				var posY float64

				// shoot up or down
				if aimY != 0 {
					moveY = -spellBulletSpeed
					posX = r.x
					posY = r.y - 8
					if aimY == 1 {
						moveY = spellBulletSpeed
						posY = r.y + 16
					}
				} else {
					moveX = spellBulletSpeed
					posX = r.x + 8
					if r.isFlip {
						moveX = moveX * -1
						posX = r.x - 8
					}
					posY = r.y - 2
					if r.isCrouch {
						posY = r.y + 6
					}
				}
				spellObj := NewSpellObject(game, posX, posY, moveX, moveY)
				game.AddSpellObject(spellObj)
			}
		}
	}

	game.debug.DrawBox(color.Black, r.x, r.y, common.TileSize, common.TileSize)
}

func (r *Player) Draw(camera common.Camera) {
	if r.postDamageTimer > 0 || r.takeDamageTimer > 0 {
		if math.Mod(r.postDamageTimer, 0.16) > 0.08 {
			return
		}
	}
	op := &ebiten.DrawImageOptions{}
	if r.isFlip {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(r.drawSizex), 0)
	}

	op.GeoM.Translate(r.x-r.drawOffsetX, r.y-r.drawOffsetY)
	op.GeoM.Scale(common.Scale, common.Scale)

	camera.DrawImage(r.animations[r.currentAnimation].GetCurrentFrame(), op)
}

func (r *Player) GetPos() (float64, float64) {
	return r.x, r.y
}

func (r *Player) TakeDamage(game *Game) {
	// already busy taking damage
	if r.takeDamageTimer > 0 {
		return
	}
	// a bit of iframes after damage
	if r.postDamageTimer > 0 {
		return
	}
	r.Health -= 1
	if r.Health > 0 {
		r.takeDamageTimer = takeDamageTime
		r.ForceJump()
	} else {
		r.deathTimer = playerDeathTime
		r.state = dyingState
	}
}

func (r *Player) ForceJump() {
	r.alreadyAbortedJump = true
	r.velocityY = (2 * forcedJumpHeight) / (standardJumpTime)
}

func (r *Player) AddHealth(amount int) {
	r.Health += amount
	if r.Health > r.MaxHealth {
		r.Health = r.MaxHealth
	}
}

func (r *Player) AddSpell(spell string) {
	_, ok := r.spells[spell]
	if !ok {
		r.spells[spell] = true
	}
	if len(r.spells) == 1 {
		r.currentSpell = spell
	}
}
