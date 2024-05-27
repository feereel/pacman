package logic

import (
	"github.com/feereel/pacman/internal/gamemap"
	"github.com/feereel/pacman/internal/models"
)

func ProcessFrame(players []*models.Player, gameMap *gamemap.GameMap) {
	for _, player := range players {
		if CanMoveForward(player, players, gameMap) {
			gameMap.SetCell(player.Position, gamemap.Empty)
			player.MoveForward()
		}
		EatFood(player, gameMap)

		gameMap.SetCell(player.Position, gamemap.Player)
	}
}

func CanMoveForward(player *models.Player, players []*models.Player, gameMap *gamemap.GameMap) bool {
	newPosition := player.Position.Add(player.Direction)
	var canMove = newPosition.InBound(0, gameMap.Width, 0, gameMap.Height) &&
		gameMap.GetCell(newPosition) != gamemap.Wall
	if !canMove {
		return false
	}

	for _, p := range players {
		if newPosition == p.Position && p.Name != player.Name {
			return false
		}
	}
	return true
}

func EatFood(player *models.Player, gameMap *gamemap.GameMap) {
	if gameMap.GetCell(player.Position) != gamemap.Food {
		return
	}
	gameMap.SetCell(player.Position, gamemap.Empty)
	gameMap.FoodCount--
	player.Score++
}
