package main

import (
	"experiments/experiments/RPG/game"
	"experiments/experiments/RPG/ui2d"
)

func main() {
	game.Run(&ui2d.UI2d{})
}
