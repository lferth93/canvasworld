package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	game := newGame(500, 500, 60)
	game.loadWorld("w1.json")
	for !rl.WindowShouldClose() {
		game.readInput()
		game.update()
		game.draw()
	}
	rl.CloseWindow()

}
