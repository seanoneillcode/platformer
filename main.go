package main

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"platformer/common"
)

func main() {
	runner := NewRunner()

	ebiten.SetWindowSize(common.ScreenWidth*common.Scale, common.ScreenHeight*common.Scale)
	ebiten.SetWindowTitle("Platform Game")
	err := ebiten.RunGame(runner)
	if err != nil {
		if errors.Is(err, common.NormalEscapeError) {
			log.Println("exiting normally")
		} else {
			log.Fatal(err)
		}
	}
}
