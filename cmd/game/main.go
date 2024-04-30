package main

import (
	"fmt"

	"github.com/feereel/pacman/internal/gamemap"
)

func main() {
	gameMap, _ := gamemap.NewSymmetricGameMap(10, 10, 1)
	fmt.Println(gameMap)
}
