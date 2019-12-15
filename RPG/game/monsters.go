package game

import "fmt"

type Monster struct {
	Character
}

func NewRat(pos Position) *Monster {
	return &Monster{Character{
		Entity: Entity{
			Position: pos,
			Name:     "Rat",
			Rune:     'R',
		},
		Hitpoints:    500,
		Strength:     0,
		Speed:        1.5,
		ActionPoints: 0.0,
	}}
}

func NewSpider(pos Position) *Monster {
	return &Monster{Character{
		Entity: Entity{
			Position: pos,
			Name:     "Spider",
			Rune:     'S',
		},
		Hitpoints:    1000,
		Strength:     0,
		Speed:        1.0,
		ActionPoints: 0.0,
	}}
}

func (m *Monster) Update(level *Level) {
	m.ActionPoints += m.Speed
	var (
		playerPos = level.Player.Position
		pos       = level.astar(m.Position, playerPos)
	)

	var movIndex = 1
	for i := 0; i < int(m.ActionPoints); i++ {
		// Most be > 1 because 1st position is the monsters current
		if movIndex < len(pos) {
			m.Move(pos[movIndex], level)
			movIndex++
			m.ActionPoints--
		}
	}
}

func (m *Monster) Move(pos Position, level *Level) {
	if _, ok := level.Monsters[pos]; !ok && pos != level.Player.Position {
		delete(level.Monsters, m.Position)
		level.Monsters[pos] = m
		m.Position = pos
	} else {
		level.AddEvent(fmt.Sprintf("%s Attacks %d Player !", m.Name, m.Strength))
		Attack(m, level.Player)
		if m.Hitpoints <= 0 {
			delete(level.Monsters, m.Position)
		}
	}
}
