package main

import (
	"experiments/experiments/RPG/game"
	"experiments/experiments/RPG/ui2d"
)

func main() {
	g := game.NewGame(1, "C:/Users/xpoc_/go/src/experiments/experiments/RPG/game/maps/level_1.map")

	go g.Run()

	ui2d.NewUI(g.InputChan, g.LevelChans[0]).Run()
}
