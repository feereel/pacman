package main

import (
	"fmt"

	"github.com/feereel/pacman/internal/gamemap"
	"github.com/feereel/pacman/internal/logic"
	"github.com/feereel/pacman/internal/models"
)

func main() {
	gameMap, _ := gamemap.NewGameMap(10, 8, 0.5)
	fmt.Println(gameMap)

	player := models.NewPlayer("Alexander", 0, 0, models.MoveRight)
	for {
		var i string
		fmt.Scan(&i)
		switch i {
		case "w":
			player.SetDirection(models.MoveUp)
		case "d":
			player.SetDirection(models.MoveRight)
		case "a":
			player.SetDirection(models.MoveLeft)
		case "s":
			player.SetDirection(models.MoveDown)
		}

		logic.ProcessFrame([]*models.Player{&player}, &gameMap)

		fmt.Println(gameMap)
		fmt.Println(player)
	}

}
