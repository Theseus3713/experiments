package game

import (
	"bufio"
	"fmt"
	"os"
)

type GameUI interface {
	Draw(*Level)
	GetInput() *Input
}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	Quit
)

type Input struct {
	Type InputType
}
type Title rune

const (
	StoneWall Title = '#'
	DirtFloor Title = '.'
	CloseDoor Title = '|'
	OpenDoor  Title = '/'
	Blank     Title = 0
	Pending   Title = -1
)

type Level struct {
	Map [][]Title
	Player
}

type Player struct {
	Entity
}

type Entity struct {
	X, Y int
}

func loadLevelFromFile(fileName string) *Level {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var (
		scanner    = bufio.NewScanner(file)
		levelLines = make([]string, 0)
		longestRaw = 0
		index      = 0
	)
	for scanner.Scan() {
		levelLines = append(levelLines, scanner.Text())
		if count := len(levelLines[index]); count > longestRaw {
			longestRaw = count
		}
		index++
	}

	level := &Level{}
	level.Map = make([][]Title, len(levelLines))
	for i := range level.Map {
		level.Map[i] = make([]Title, longestRaw)
	}
	var lenMap = len(level.Map)
	for y := 0; y < lenMap; y++ {
		var line = levelLines[y]
		for x, c := range line {
			var t Title
			switch c {
			case ' ', '\t', '\n', '\r':
				t = Blank
			case '#':
				t = StoneWall
			case '|':
				t = CloseDoor
			case '/':
				t = OpenDoor
			case '.':
				t = DirtFloor
			case 'P':
				level.Player.X = x
				level.Player.Y = y
				t = Pending
			default:
				panic(fmt.Sprintf(`Invalid character '%d' in map`, c))
			}
			level.Map[y][x] = t
		}
	}
	for y, row := range level.Map {
		for x, tile := range row {
			if tile == Pending {
			SearchLoop:
				for searchX := x - 1; searchX <= x+1; searchX++ {
					for searchY := x - 1; searchY <= x+1; searchY++ {
						var searchTile = level.Map[searchX][searchY]
						switch searchTile {
						case DirtFloor:
							level.Map[y][x] = DirtFloor
							break SearchLoop
						}
					}
				}
			}
		}
	}

	return level

}

func canWalk(level *Level, x, y int) bool {
	switch level.Map[y][x] {
	case StoneWall, CloseDoor, Blank:
		return false
	default:
		return true
	}
}

func checkDoor(level *Level, x, y int) {
	if level.Map[y][x] == CloseDoor {
		level.Map[y][x] = OpenDoor
	}
}

func handleInput(level *Level, input *Input) {
	var player = level.Player
	switch input.Type {
	case Up:
		if canWalk(level, player.X, player.Y-1) {
			level.Player.Y--
		} else {
			checkDoor(level, player.X, player.Y-1)
		}
	case Down:
		if canWalk(level, player.X, player.Y+1) {
			level.Player.Y++
		} else {
			checkDoor(level, player.X, player.Y+1)
		}
	case Left:
		if canWalk(level, player.X-1, player.Y) {
			level.Player.X--
		} else {
			checkDoor(level, player.X-1, player.Y)
		}
	case Right:
		if canWalk(level, player.X+1, player.Y) {
			level.Player.X++
		} else {
			checkDoor(level, player.X+1, player.Y)
		}
	}
}

func Run(ui GameUI) {
	var level = loadLevelFromFile("C:/Users/xpoc_/go/src/experiments/experiments/RPG/game/maps/level_1.map")
	for {
		ui.Draw(level)
		input := ui.GetInput()

		if input.Type == Quit {
			return
		}
		handleInput(level, input)
	}
}
