package models

import (
	"net"

	"github.com/feereel/pacman/internal/utility"
)

type MoveDirection byte

const (
	MoveUp    MoveDirection = 0
	MoveRight MoveDirection = 1
	MoveDown  MoveDirection = 2
	MoveLeft  MoveDirection = 3
)

var DirToVec = map[MoveDirection]utility.Vector2D[int]{
	MoveUp:    {X: 0, Y: -1},
	MoveRight: {X: 1, Y: 0},
	MoveDown:  {X: 0, Y: 1},
	MoveLeft:  {X: -1, Y: 0},
}

var VecToDir = map[utility.Vector2D[int]]MoveDirection{
	{X: 0, Y: -1}: MoveUp,
	{X: 1, Y: 0}:  MoveRight,
	{X: 0, Y: 1}:  MoveDown,
	{X: -1, Y: 0}: MoveLeft,
}

type Player struct {
	Name       string
	Score      int
	Controlled bool
	Conn       net.Conn

	Position  utility.Vector2D[int]
	Direction utility.Vector2D[int]
}

func NewPlayer(name string, startX int, startY int, direction MoveDirection) Player {
	player := Player{
		Name:       name,
		Score:      0,
		Position:   utility.Vector2D[int]{X: startX, Y: startY},
		Direction:  DirToVec[direction],
		Controlled: false,
	}

	return player
}

func (player *Player) MoveForward() {
	player.Position.X += player.Direction.X
	player.Position.Y += player.Direction.Y
}

func (player *Player) SetDirection(direction MoveDirection) {
	player.Direction = DirToVec[direction]
}

func NameInPlayers(name string, players []Player) (int, bool) {
	for i, p := range players {
		if p.Name == name {
			return i, true
		}
	}
	return -1, false
}

func NameInLPlayers(name string, players []*Player) (int, bool) {
	for i, p := range players {
		if p.Name == name {
			return i, true
		}
	}
	return -1, false
}
