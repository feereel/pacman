package network

import (
	"net"

	"github.com/feereel/pacman/internal/gamemap"
)

func SerializeGameMap(gameMap gamemap.GameMap) ([]byte, error) {
	data := make([]byte, gameMap.Height*gameMap.Width)
	for y := 0; y < gameMap.Height; y++ {
		for x := 0; x < gameMap.Width; x++ {
			i := y*gameMap.Width + x
			data[i] = byte(gameMap.Grid[y][x])
		}
	}
	return data, nil
}

func DeserializeGameMap(conn net.Conn, width int, height int) (gamemap.GameMap, error) {
	data := make([]byte, width*height)
	n, err := conn.Read(data)
	if n != width*height || err != nil {
		return gamemap.GameMap{}, err
	}

	grid := make([][]gamemap.Cell, height)
	for i := range grid {
		grid[i] = make([]gamemap.Cell, width)
	}
	gameMap := gamemap.GameMap{Width: width, Height: height, Grid: grid}

	for y := 0; y < gameMap.Height; y++ {
		for x := 0; x < gameMap.Width; x++ {
			cell := gamemap.Cell(data[y*gameMap.Width+x])
			gameMap.Grid[y][x] = cell
			if cell == gamemap.Food {
				gameMap.FoodCount++
			}
		}
	}
	return gameMap, nil
}
