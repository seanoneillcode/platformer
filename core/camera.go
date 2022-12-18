package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"platformer/common"
)

type Camera struct {
	x      float64
	y      float64
	buffer *ebiten.Image
	target CameraTarget
}

type CameraTarget interface {
	GetPos() (float64, float64)
}

func NewCamera() *Camera {
	return &Camera{
		buffer: ebiten.NewImage(common.ScreenWidth*common.Scale, common.ScreenHeight*common.Scale),
	}
}

func (c *Camera) Update(delta float64, game *Game) {
	c.buffer.Clear()
	if c.target != nil {
		tx, ty := c.target.GetPos()
		c.x = tx - (common.ScreenWidth / 2) + (common.TileSize / 2)
		c.y = ty - (common.ScreenHeight / 2) + (common.TileSize / 2)
	}
	if c.x < 0 {
		c.x = 0
	}
	if c.y < 0 {
		c.y = 0
	}
	maxWidth := float64((game.level.tiledGrid.GroundLayer.Width * common.TileSize) - common.ScreenWidth)
	if c.x > maxWidth {
		c.x = maxWidth
	}
	maxHeight := float64((game.level.tiledGrid.GroundLayer.Height * common.TileSize) - common.ScreenHeight)
	if c.y > maxHeight {
		c.y = maxHeight
	}
}

func (c *Camera) DrawBuffer(screen *ebiten.Image) {
	ops := &ebiten.DrawImageOptions{}
	screen.DrawImage(c.buffer, ops)
}

func (c *Camera) DrawImage(img *ebiten.Image, options *ebiten.DrawImageOptions) {
	options.GeoM.Translate(-c.x*common.Scale, -c.y*common.Scale)
	c.buffer.DrawImage(img, options)
}

func (c *Camera) Target(target CameraTarget) {
	c.target = target
}

func (c *Camera) GetPos() (float64, float64) {
	return c.x, c.y
}
