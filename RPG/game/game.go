package game

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Game struct {
	LevelChans []chan *Level
	InputChan  chan *Input
	Level      *Level
}

func NewGame(numWindows int, path string) *Game {
	levelChans := make([]chan *Level, numWindows)
	for i := range levelChans {
		levelChans[i] = make(chan *Level)
	}
	inputChan := make(chan *Input)
	return &Game{levelChans, inputChan, loadLevelFromFile(path)}
}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	QuitGame
	CloseWindow
	Search //temporary
)

type Input struct {
	Type         InputType
	LevelChannel chan *Level
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
	Debug map[Position]bool
}

type Player struct {
	Entity
}

type Position struct {
	X, Y int
}

type Entity struct {
	Position
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

func canWalk(level *Level, pos Position) bool {
	switch level.Map[pos.Y][pos.X] {
	case StoneWall, CloseDoor, Blank:
		return false
	default:
		return true
	}
}

func checkDoor(level *Level, pos Position) {
	if level.Map[pos.Y][pos.X] == CloseDoor {
		level.Map[pos.Y][pos.X] = OpenDoor
	}
}

func (game *Game) handleInput(input *Input) {
	var player = game.Level.Player
	switch input.Type {
	case Up:
		if canWalk(game.Level, Position{player.X, player.Y - 1}) {
			game.Level.Player.Y--
		} else {
			checkDoor(game.Level, Position{player.X, player.Y - 1})
		}
	case Down:
		if canWalk(game.Level, Position{player.X, player.Y + 1}) {
			game.Level.Player.Y++
		} else {
			checkDoor(game.Level, Position{player.X, player.Y + 1})
		}
	case Left:
		if canWalk(game.Level, Position{player.X - 1, player.Y}) {
			game.Level.Player.X--
		} else {
			checkDoor(game.Level, Position{player.X - 1, player.Y})
		}
	case Right:
		if canWalk(game.Level, Position{player.X + 1, player.Y}) {
			game.Level.Player.X++
		} else {
			checkDoor(game.Level, Position{player.X + 1, player.Y})
		}
	case Search:
		//bfs(ui, level, player.Position)
		game.astar(player.Position, Position{3, 2})
	case CloseWindow:
		close(input.LevelChannel)
		chanIndex := 0
		for i, c := range game.LevelChans {
			if input.LevelChannel == c {
				chanIndex = i
				break
			}
		}
		game.LevelChans = append(game.LevelChans[:chanIndex], game.LevelChans[chanIndex+1:]...)

	}
}

func getNeighbors(level *Level, pos Position) []Position {
	var (
		neighbors = make([]Position, 0, 4)
		left      = Position{pos.X - 1, pos.Y}
		right     = Position{pos.X + 1, pos.Y}
		up        = Position{pos.X, pos.Y - 1}
		down      = Position{pos.X, pos.Y + 1}
	)
	if canWalk(level, right) {
		neighbors = append(neighbors, right)
	}
	if canWalk(level, left) {
		neighbors = append(neighbors, left)
	}
	if canWalk(level, up) {
		neighbors = append(neighbors, up)
	}
	if canWalk(level, down) {
		neighbors = append(neighbors, down)
	}
	return neighbors
}

func (game *Game) bfs(start Position) {
	var frontier = make([]Position, 0, 8)
	frontier = append(frontier, start)
	var visited = make(map[Position]bool)
	visited[start] = true
	game.Level.Debug = visited

	for len(frontier) > 0 {
		var current = frontier[0]
		frontier = frontier[1:]
		for _, next := range getNeighbors(game.Level, current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
			}
		}
	}
}

func (game *Game) astar(start, goal Position) []Position {
	var frontier = make(pQueue, 0, 8)
	frontier = frontier.push(start, 1)
	//frontier = append(frontier, priorityPosition{start, 1})
	var cameFrom = make(map[Position]Position)
	cameFrom[start] = start
	var costSoFor = make(map[Position]int)
	costSoFor[start] = 0

	game.Level.Debug = make(map[Position]bool)
	var current Position
	for len(frontier) > 0 {
		frontier, current = frontier.pop()
		if current == goal {
			var path = make([]Position, 0)
			var pos = current
			for pos != start {
				path = append(path, pos)
				pos = cameFrom[pos]
			}
			path = append(path, pos)

			// Инверсия [3,2,1] = [1,2,3]
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}

			for _, pos := range path {
				game.Level.Debug[pos] = true
			}
			return path
		}
		for _, next := range getNeighbors(game.Level, current) {
			var newCost = costSoFor[current] + 1 // always 1 for now
			if _, ok := costSoFor[next]; !ok || newCost < costSoFor[next] {
				costSoFor[next] = newCost
				var (
					xDist    = int(math.Abs(float64(goal.X - next.X)))
					yDist    = int(math.Abs(float64(goal.Y - next.Y)))
					priority = newCost + xDist + yDist
				)
				frontier = frontier.push(next, priority)
				cameFrom[next] = current

			}
		}
	}
	return nil
}

func (game *Game) Run() {
	fmt.Println("Starting...")

	for _, lChan := range game.LevelChans {
		lChan <- game.Level
	}
	for input := range game.InputChan {
		if input.Type == QuitGame {
			return
		}
		if input.Type == CloseWindow {
			// TODO 	windows.Close() fd Handle in ui
			return
		}
		game.handleInput(input)
		for _, lChan := range game.LevelChans {
			lChan <- game.Level
		}
	}
}
