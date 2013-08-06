package main

import (
	"github.com/diasf/pongo/game"
)

func main() {
	game := game.NewPongoGame(800, 600)
	game.Start()
}
