package logic

import (
	"fmt"

	"github.com/feereel/pacman/internal/gamemap"
	"github.com/feereel/pacman/internal/models"
)

func ProcessFrame(players []*models.Player, gameMap *gamemap.GameMap) {
	for _, player := range players {
		if CanMoveForward(player, gameMap) {
			gameMap.SetCell(player.Position, gamemap.Empty)
			player.MoveForward()
		}
		EatFood(player, gameMap)

		gameMap.SetCell(player.Position, gamemap.Player)
	}
	fmt.Println(players[0])
}

func CanMoveForward(player *models.Player, gameMap *gamemap.GameMap) bool {
	newPosition := player.Position.Add(player.Direction)
	return newPosition.InBound(0, gameMap.Width, 0, gameMap.Height) && gameMap.GetCell(newPosition) != gamemap.Wall
}

func EatFood(player *models.Player, gameMap *gamemap.GameMap) {
	if gameMap.GetCell(player.Position) != gamemap.Food {
		return
	}
	gameMap.SetCell(player.Position, gamemap.Empty)
	player.Score++
}
