package main

import (
	"github.com/pongo/game"
)

func main() {
	game := game.NewPongoGame(800, 600)
	game.Start()
}
