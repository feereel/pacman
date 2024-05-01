package client

import (
	"fmt"
	"os"

	"github.com/feereel/pacman/internal/gamemap"
	"github.com/feereel/pacman/internal/logic"
	"github.com/feereel/pacman/internal/models"
	"github.com/feereel/pacman/internal/terminal"
	termbox "github.com/nsf/termbox-go"
)

func Run() int {
	gameMap, err := gamemap.NewSymmetricGameMap(20, 15, 0.2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	player := models.NewPlayer("Alexander", 0, 0, models.MoveRight)
	player.Controlled = true

	term, err := terminal.NewTerminal(&gameMap, []*models.Player{&player, &player})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	defer term.Close()

	for {
		term.Render()
		e := term.Poll()
		if e.Ch == 'q' || e.Key == termbox.KeyCtrlC {
			return 0
		}
		handleKeypress(&player, e)
		logic.ProcessFrame([]*models.Player{&player}, &gameMap)
	}
}

func handleKeypress(player *models.Player, e termbox.Event) {
	if e.Ch == 0 {
		switch e.Key {
		case termbox.KeyArrowUp:
			player.SetDirection(models.MoveUp)
		case termbox.KeyArrowDown:
			player.SetDirection(models.MoveDown)
		case termbox.KeyArrowLeft:
			player.SetDirection(models.MoveLeft)
		case termbox.KeyArrowRight:
			player.SetDirection(models.MoveRight)
		}
	}
}
