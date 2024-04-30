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

var CellVisualSymbol = map[Cell]string{
	Empty:  "_",
	Player: "@",
	Food:   ".",
	Wall:   "#",
}

type GameMap struct {
	Grid   [][]Cell
	Width  int32
	Height int32
}

func NewSymmetricGameMap(Width int32, Height int32, occupancyPercantage float32) (GameMap, error) {
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
		Grid:   grid,
		Width:  symWidth,
		Height: symHeight,
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

func NewGameMap(Width int32, Height int32, occupancyPercantage float32) (GameMap, error) {
	if Width <= 0 || Height <= 0 {
		return GameMap{}, errors.New("width or heigth is less than 1")
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

	var wallsCount int = int(float32(gameMap.Height*gameMap.Width-gameMap.perimeter()) * occupancyPercantage)

	gameMap.GenerateMap(wallsCount)
	gameMap.Grid[0][0] = Player

	return gameMap, nil
}

func (gameMap *GameMap) GenerateMap(wallsCount int) {
	var currentObstacleCount int = 0
	for i := 0; i < wallsCount; i++ {
		x := rand.Int31()%(gameMap.Width-2) + 1
		y := rand.Int31()%(gameMap.Height-2) + 1

		if gameMap.Grid[y][x] != Empty {
			continue
		}

		gameMap.Grid[y][x] = Wall
		currentObstacleCount++

		if !gameMap.IsFullyAccessible(int32(currentObstacleCount)) {
			gameMap.Grid[y][x] = Empty
			currentObstacleCount--
		}
	}

	for y := 0; y < int(gameMap.Height); y++ {
		for x := 0; x < int(gameMap.Width); x++ {
			if gameMap.Grid[y][x] == Empty {
				gameMap.Grid[y][x] = Food
			}
		}
	}

}

func (gameMap GameMap) IsFullyAccessible(currentObstacleCount int32) bool {

	mapFlags := make([][]bool, gameMap.Height)
	for i := range mapFlags {
		mapFlags[i] = make([]bool, gameMap.Width)
	}

	queue := make([]utility.Vector2D[int32], 0)
	var i int32
	for i = 0; i < gameMap.Width; i++ {
		queue = append(queue, utility.Vector2D[int32]{X: i, Y: 0})
		queue = append(queue, utility.Vector2D[int32]{X: i, Y: gameMap.Height - 1})
		mapFlags[0][i] = true
		mapFlags[gameMap.Height-1][i] = true
	}
	for i = 0; i < gameMap.Height; i++ {
		queue = append(queue, utility.Vector2D[int32]{X: 0, Y: i})
		queue = append(queue, utility.Vector2D[int32]{X: gameMap.Width - 1, Y: i})
		mapFlags[i][0] = true
		mapFlags[i][gameMap.Width-1] = true
	}

	var accesibleCellCount int32 = gameMap.perimeter()

	for len(queue) > 0 {
		cell := queue[0]
		queue = queue[1:]

		var x, y int32
		for x = -1; x < 1; x++ {
			for y = -1; y < 1; y++ {
				neighbourX := cell.X + x
				neighbourY := cell.Y + y
				if x != 0 && y != 0 {
					continue
				}
				if neighbourX < 0 || neighbourX >= gameMap.Width || neighbourY < 0 || neighbourY >= gameMap.Height {
					continue
				}
				if mapFlags[neighbourY][neighbourX] || gameMap.Grid[neighbourY][neighbourX] != Empty {
					continue
				}
				mapFlags[neighbourY][neighbourX] = true
				queue = append(queue, utility.Vector2D[int32]{X: neighbourX, Y: neighbourY})
				accesibleCellCount++
			}
		}
	}

	var targetAccessibleCellCount int32 = gameMap.Width*gameMap.Height - currentObstacleCount
	return targetAccessibleCellCount == accesibleCellCount
}

func (gameMap GameMap) perimeter() int32 {
	return gameMap.Width*2 + gameMap.Height*2 - 4
}

func (gameMap GameMap) String() string {

	var result string

	for y := 0; y < int(gameMap.Height); y++ {
		for x := 0; x < int(gameMap.Width); x++ {
			result += CellVisualSymbol[gameMap.Grid[y][x]] + " "
		}
		result += "\n"
	}
	return result
}
