package models

import (
	"github.com/feereel/pacman/internal/utility"
)

type MoveDirection uint32

const (
	MoveUp    MoveDirection = 0
	MoveRight MoveDirection = 1
	MoveDown  MoveDirection = 2
	MoveLeft  MoveDirection = 3
)

var MapDirrections = map[MoveDirection]utility.Vector2D[int]{
	MoveUp:    {X: 0, Y: -1},
	MoveRight: {X: 1, Y: 0},
	MoveDown:  {X: 0, Y: 1},
	MoveLeft:  {X: -1, Y: 0},
}

type Player struct {
	Name  string
	Score int

	Position  utility.Vector2D[int]
	Direction utility.Vector2D[int]
}

func NewPlayer(name string, startX int, startY int, direction MoveDirection) Player {
	player := Player{
		Name:      name,
		Score:     0,
		Position:  utility.Vector2D[int]{X: startX, Y: startY},
		Direction: MapDirrections[direction],
	}

	return player
}

func (player *Player) MoveForward() {
	player.Position.X += player.Direction.X
	player.Position.Y += player.Direction.Y
}

func (player *Player) SetDirection(direction MoveDirection) {
	player.Direction = MapDirrections[direction]
}
