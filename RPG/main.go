package main

import (
	"experiments/experiments/RPG/game"
	"experiments/experiments/RPG/ui2d"
	"runtime"
)

func main() {
	g := game.NewGame(1, "C:/Users/xpoc_/go/src/experiments/experiments/RPG/game/maps/level_1.map")
	for i := 0; i < 1; i++ {
		go func(i int) {
			runtime.LockOSThread()
			ui2d.NewUI(g.InputChan, g.LevelChans[i]).Run()
		}(i)
	}
	g.Run()
}
