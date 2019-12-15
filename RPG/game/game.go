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
	Map           [][]Title
	Player        *Player
	Monsters      map[Position]*Monster
	Debug         map[Position]bool
	Events        []string
	EventPosition int
}

type Player struct {
	Character
}

type Attackable interface {
	GetActionPints() float64
	SetActionPints(float64)
	GetHitpoints() int
	SetHitpoints(int)
	GetAttackPower() int
}

func (c *Character) GetActionPints() float64 {
	return c.ActionPoints
}

func (c *Character) SetActionPints(ap float64) {
	c.ActionPoints = ap
}

func (c *Character) GetHitpoints() int {
	return c.Hitpoints
}

func (c *Character) SetHitpoints(hp int) {
	c.Hitpoints = hp
}

func (c *Character) GetAttackPower() int {
	return c.Strength
}

func Attack(a1, a2 Attackable) {
	a1.SetActionPints(a1.GetActionPints() - 1)
	a2.SetHitpoints(a2.GetHitpoints() - a1.GetAttackPower())
	if a2.GetHitpoints() > 0 {
		a2.SetActionPints(a2.GetActionPints() - 1)
		a1.SetHitpoints(a1.GetHitpoints() - a2.GetAttackPower())
	}
	fmt.Println("PlayerAttackMonster")
}

func (level *Level) AddEvent(event string) {
	level.Events[level.EventPosition] = event
	level.EventPosition++
	if level.EventPosition == len(level.Events) {
		level.EventPosition = 0
	}
}

type Position struct {
	X, Y int
}

type Entity struct {
	Position
	Name string
	Rune rune
}

type Character struct {
	Entity
	Hitpoints    int
	Strength     int
	Speed        float64
	ActionPoints float64
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

	level := &Level{
		Player: &Player{
			Character{
				Entity: Entity{
					Name: "GoMen",
					Rune: '@',
				},
				Hitpoints:    20,
				Strength:     20,
				Speed:        1.0,
				ActionPoints: 0,
			}},
		Events: make([]string, 10),
	}

	level.Map = make([][]Title, len(levelLines))
	level.Monsters = make(map[Position]*Monster)
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
			case '@':
				level.Player.X = x
				level.Player.Y = y
				t = Pending
			case 'R':
				level.Monsters[Position{x, y}] = NewRat(Position{x, y})
				t = Pending
			case 'S':
				level.Monsters[Position{x, y}] = NewSpider(Position{x, y})
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
				level.Map[y][x] = level.bfsFloor(Position{x, y})
			}
		}
	}

	return level

}

func inRange(level *Level, pos Position) bool {
	return pos.X < len(level.Map[0]) && pos.Y < len(level.Map) && pos.X >= 0 && pos.Y >= 0
}

func canWalk(level *Level, pos Position) bool {
	if inRange(level, pos) {
		switch level.Map[pos.Y][pos.X] {
		case StoneWall, CloseDoor, Blank:
			return false
		default:
			return true
		}
	}
	return false
}

func checkDoor(level *Level, pos Position) {
	if level.Map[pos.Y][pos.X] == CloseDoor {
		level.Map[pos.Y][pos.X] = OpenDoor
	}
}

func (p *Player) Move(pos Position, level *Level) {
	if monsters, ok := level.Monsters[pos]; !ok {
		p.Position = pos
	} else {
		Attack(level.Player, monsters)
		level.AddEvent("Player Attacked Monster")
		if monsters.Hitpoints <= 0 {
			delete(level.Monsters, monsters.Position)
		}
		if level.Player.Hitpoints <= 0 {
			fmt.Println("YOU DIED")
			panic("YOU DIED")
		}
	}
}

func (game *Game) handleInput(input *Input) {
	var (
		level  = game.Level
		player = level.Player
	)
	switch input.Type {
	case Up:
		newPos := Position{player.X, player.Y - 1}
		if canWalk(game.Level, newPos) {
			level.Player.Move(newPos, level)
		} else {
			checkDoor(game.Level, newPos)
		}
	case Down:
		newPos := Position{player.X, player.Y + 1}
		if canWalk(game.Level, newPos) {
			player.Move(newPos, level)
		} else {
			checkDoor(game.Level, newPos)
		}
	case Left:
		newPos := Position{player.X - 1, player.Y}
		if canWalk(game.Level, newPos) {
			player.Move(newPos, level)
		} else {
			checkDoor(game.Level, newPos)
		}
	case Right:
		newPos := Position{player.X + 1, player.Y}
		if canWalk(game.Level, newPos) {
			player.Move(newPos, level)
		} else {
			checkDoor(game.Level, newPos)
		}
	case Search:
		//bfs(ui, level, player.Position)
		level.astar(player.Position, Position{3, 2})
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

func (level *Level) bfsFloor(start Position) Title {
	var frontier = make([]Position, 0, 8)
	frontier = append(frontier, start)
	var visited = make(map[Position]bool)
	visited[start] = true
	level.Debug = visited

	for len(frontier) > 0 {
		var current = frontier[0]

		var currentTile = level.Map[current.Y][current.X]
		switch currentTile {
		case DirtFloor:
			return DirtFloor
		default:
		}

		frontier = frontier[1:]
		for _, next := range getNeighbors(level, current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
			}
		}
	}
	return DirtFloor
}

func (level *Level) astar(start, goal Position) []Position {
	var frontier = make(pQueue, 0, 8)
	frontier = frontier.push(start, 1)
	var cameFrom = make(map[Position]Position)
	cameFrom[start] = start
	var costSoFor = make(map[Position]int)
	costSoFor[start] = 0

	//level.Debug = make(map[Position]bool)
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

			//for _, pos := range path {
			//	level.Debug[pos] = true
			//}
			return path
		}
		for _, next := range getNeighbors(level, current) {
			var newCost = costSoFor[current] + 1 // always 1 for now
			if _, ok := costSoFor[next]; !ok || newCost < costSoFor[next] {
				costSoFor[next] = newCost
				var (
					xDist    = int(math.Abs(float64(goal.X - next.X)))
					yDist    = int(math.Abs(float64(goal.Y - next.Y)))
					priority = newCost + xDist + yDist
				)
				frontier = frontier.push(next, priority)
				//level.Debug[next] = true
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

		for _, monster := range game.Level.Monsters {
			monster.Update(game.Level)
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
