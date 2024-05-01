package gamemap

import (
	"errors"
	"math/rand"

	"github.com/feereel/pacman/internal/utility"
)

type Cell uint8

const (
	Empty  Cell = 0x0
	Player Cell = 0x22
	Food   Cell = 0xAA
	Wall   Cell = 0xFF
)

type GameMap struct {
	Grid      [][]Cell
	Width     int
	Height    int
	FoodCount int
}

func NewSymmetricGameMap(Width int, Height int, occupancyPercantage float32) (GameMap, error) {
	gameMap, err := NewGameMap(Width, Height, occupancyPercantage)
	if err != nil {
		return gameMap, err
	}

	symWidth := Width * 2
	symHeight := Height * 2

	grid := make([][]Cell, symHeight)
	for i := range grid {
		grid[i] = make([]Cell, symWidth)
	}

	symGameMap := GameMap{
		Grid:      grid,
		Width:     symWidth,
		Height:    symHeight,
		FoodCount: gameMap.FoodCount * 2,
	}

	for y := 0; y < int(Height); y++ {
		for x := 0; x < int(Width); x++ {
			symGameMap.Grid[y][x] = gameMap.Grid[y][x]
			symGameMap.Grid[y][int(symWidth)-x-1] = gameMap.Grid[y][x]
			symGameMap.Grid[int(symHeight)-y-1][x] = gameMap.Grid[y][x]
			symGameMap.Grid[int(symHeight)-y-1][int(symWidth)-x-1] = gameMap.Grid[y][x]
		}
	}

	return symGameMap, nil
}

func NewGameMap(Width int, Height int, occupancyPercantage float32) (GameMap, error) {
	if Width < 5 || Height < 5 {
		return GameMap{}, errors.New("width or heigth is less than 5")
	}

	grid := make([][]Cell, Height)
	for i := range grid {
		grid[i] = make([]Cell, Width)
	}

	gameMap := GameMap{
		Grid:   grid,
		Width:  Width,
		Height: Height,
	}

	var wallsCount int = int(float32(gameMap.Height*gameMap.Width-4) * occupancyPercantage)

	gameMap.Grid[0][0] = Player
	gameMap.Grid[gameMap.Height-1][gameMap.Width-1] = Food
	gameMap.GenerateMap(wallsCount)

	return gameMap, nil
}

func (gameMap *GameMap) GenerateMap(wallsCount int) {
	var currentObstacleCount int = 0
	for i := 0; i < wallsCount; i++ {
		x := rand.Int() % gameMap.Width
		y := rand.Int() % gameMap.Height

		if gameMap.Grid[y][x] != Empty {
			continue
		}

		gameMap.Grid[y][x] = Wall
		currentObstacleCount++
		if !gameMap.IsFullyAccessible(int(currentObstacleCount)) {
			gameMap.Grid[y][x] = Empty
			currentObstacleCount--
		}
	}

	for y := 0; y < int(gameMap.Height); y++ {
		for x := 0; x < int(gameMap.Width); x++ {
			if gameMap.Grid[y][x] == Empty {
				gameMap.Grid[y][x] = Food
				gameMap.FoodCount++
			}
		}
	}

}

func (gameMap GameMap) IsFullyAccessible(currentObstacleCount int) bool {

	mapFlags := make([][]bool, gameMap.Height)
	for i := range mapFlags {
		mapFlags[i] = make([]bool, gameMap.Width)
	}

	queue := make([]utility.Vector2D[int], 0)
	queue = append(queue, utility.Vector2D[int]{X: 0, Y: 0})

	mapFlags[0][0] = true
	var accesibleCellCount int = 1

	for len(queue) > 0 {
		cell := queue[0]
		queue = queue[1:]

		var x, y int
		for x = -1; x < 2; x++ {
			for y = -1; y < 2; y++ {
				if x != 0 && y != 0 {
					continue
				}
				neighbour := utility.Vector2D[int]{
					X: cell.X + x,
					Y: cell.Y + y,
				}
				if !neighbour.InBound(0, gameMap.Width, 0, gameMap.Height) {
					continue
				}
				if mapFlags[neighbour.Y][neighbour.X] || gameMap.GetCell(neighbour) == Wall {
					continue
				}
				mapFlags[neighbour.Y][neighbour.X] = true
				queue = append(queue, neighbour)
				accesibleCellCount++
			}
		}
	}

	var targetAccessibleCellCount int = gameMap.Width*gameMap.Height - currentObstacleCount
	return targetAccessibleCellCount == accesibleCellCount
}

func (gameMap GameMap) GetCell(position utility.Vector2D[int]) Cell {
	return gameMap.Grid[position.Y][position.X]
}

func (gameMap GameMap) SetCell(position utility.Vector2D[int], value Cell) {
	gameMap.Grid[position.Y][position.X] = value
}
