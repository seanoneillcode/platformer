package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"platformer/common"
	"platformer/core"
	"platformer/res"
)

type Book struct {
	pageImage  *ebiten.Image
	coverImage *ebiten.Image
	x          float64
	y          float64
	text       string
	title      string
	//
	enabled bool
	visible bool
}

func NewBook(resources *res.Resources) *Book {
	return &Book{
		x:          40,
		y:          8,
		pageImage:  resources.GetImage("book-page"),
		coverImage: resources.GetImage("book-cover"),
	}
}

func (r *Book) Update(delta float64, game *core.Game) {
	if !r.enabled {
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		r.visible = false
		r.enabled = false
		game.Actions.CloseBook()
	}
}

func (r *Book) Draw(screen *ebiten.Image) {
	if !r.visible {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.x, r.y)
	op.GeoM.Scale(common.Scale, common.Scale)
	screen.DrawImage(r.pageImage, op)

	paragraphX := r.x + 5
	paragraphY := r.y + 5
	common.DrawText(screen, r.text, paragraphX, paragraphY)
}

func (r *Book) Open(title string, text string) {
	r.title = title
	r.text = text
	r.enabled = true
	r.visible = true
}
