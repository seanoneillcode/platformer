package common

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

var (
	textImage           = loadImage("common/text-source.png")
	textCharacterImages = map[rune]*ebiten.Image{}

	characterInfo = map[rune]charInfo{
		'a':  {index: 0, width: 5},
		'b':  {index: 1, width: 5},
		'c':  {index: 2, width: 5},
		'd':  {index: 3, width: 5},
		'e':  {index: 4, width: 5},
		'f':  {index: 5, width: 5},
		'g':  {index: 6, width: 5},
		'h':  {index: 7, width: 5},
		'i':  {index: 8, width: 3},
		'j':  {index: 9, width: 3},
		'k':  {index: 10, width: 4},
		'l':  {index: 11, width: 3},
		'm':  {index: 12, width: 5},
		'n':  {index: 13, width: 5},
		'o':  {index: 14, width: 5},
		'p':  {index: 15, width: 5},
		'q':  {index: 16, width: 5},
		'r':  {index: 17, width: 6},
		's':  {index: 18, width: 5},
		't':  {index: 19, width: 4},
		'u':  {index: 20, width: 5},
		'v':  {index: 21, width: 5},
		'w':  {index: 22, width: 5},
		'x':  {index: 23, width: 6},
		'y':  {index: 24, width: 5},
		'z':  {index: 25, width: 5},
		'0':  {index: 26, width: 5},
		'1':  {index: 27, width: 3},
		'2':  {index: 28, width: 5},
		'3':  {index: 29, width: 5},
		'4':  {index: 30, width: 5},
		'5':  {index: 31, width: 5},
		'6':  {index: 32, width: 5},
		'7':  {index: 33, width: 5},
		'8':  {index: 34, width: 5},
		'9':  {index: 35, width: 5},
		',':  {index: 36, width: 2},
		'.':  {index: 37, width: 1},
		'!':  {index: 38, width: 1},
		'?':  {index: 39, width: 5},
		'A':  {index: 40, width: 5},
		'B':  {index: 41, width: 5},
		'C':  {index: 42, width: 5},
		'D':  {index: 43, width: 5},
		'E':  {index: 44, width: 5},
		'F':  {index: 45, width: 5},
		'G':  {index: 46, width: 5},
		'H':  {index: 47, width: 5},
		'I':  {index: 48, width: 3},
		'J':  {index: 49, width: 5},
		'K':  {index: 50, width: 5},
		'L':  {index: 51, width: 4},
		'M':  {index: 52, width: 5},
		'N':  {index: 53, width: 5},
		'O':  {index: 54, width: 5},
		'P':  {index: 55, width: 5},
		'Q':  {index: 56, width: 5},
		'R':  {index: 57, width: 6},
		'S':  {index: 58, width: 5},
		'T':  {index: 59, width: 5},
		'U':  {index: 60, width: 5},
		'V':  {index: 61, width: 5},
		'W':  {index: 62, width: 5},
		'X':  {index: 63, width: 6},
		'Y':  {index: 64, width: 5},
		'Z':  {index: 65, width: 5},
		'\'': {index: 66, width: 2},
	}
)

type charInfo struct {
	index int
	width int
}

type Drawable interface {
	DrawImage(img *ebiten.Image, options *ebiten.DrawImageOptions)
}

func DrawText(drawable Drawable, str string, x float64, y float64) {
	drawText(drawable, str, x, y, 1)
}

func DrawTextWithAlpha(drawable Drawable, str string, x float64, y float64, alpha float64) {
	drawText(drawable, str, x, y, alpha)
}

func TextWidth(text string) int {
	x := 0
	for _, c := range text {
		if c == '\n' {
			return x / 2
		}
		if c == ' ' {
			x += spaceWidth
			continue
		}
		ci, ok := characterInfo[c]
		if !ok {
			continue
		}
		x += ci.width + 1
	}
	return x / 2
}

const (
	ch         = 12
	cw         = 10
	spaceWidth = 4
)

func drawText(drawable Drawable, str string, ox, oy float64, alpha float64) {
	x := 0
	y := 0
	for _, c := range str {
		if c == '\n' {
			x = 0
			y += ch
			continue
		}
		if c == ' ' {
			x += spaceWidth
			continue
		}
		ci, ok := characterInfo[c]
		if !ok {
			continue
		}
		s, ok := textCharacterImages[c]
		if !ok {
			sx := ci.index * cw
			rect := image.Rect(sx, 0, sx+ci.width, ch-1)
			s = textImage.SubImage(rect).(*ebiten.Image)
			textCharacterImages[c] = s
		}
		if s != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(ox*2), float64(oy*2))
			op.GeoM.Translate(float64(x), float64(y))
			op.GeoM.Scale(TextScale, TextScale)
			op.ColorM.Scale(0, 0, 0, alpha)
			drawable.DrawImage(s, op)
			x += ci.width + 1
		}
	}
}
