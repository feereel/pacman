package client

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/feereel/pacman/internal/gamemap"
	"github.com/feereel/pacman/internal/logic"
	"github.com/feereel/pacman/internal/models"
	"github.com/feereel/pacman/internal/network"
	"github.com/feereel/pacman/internal/terminal"
	"github.com/nsf/termbox-go"
)

var netclock *network.Netclock

func Run(clientName, ip string, port, mapWidth, mapHeight int) int {

	serverConn, err := net.Dial("tcp", fmt.Sprintf("%s:%v", ip, port))
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer serverConn.Close()

	err = network.SendPlayerName(serverConn, clientName)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	gameMiniMap, err := network.RecvGameMap(serverConn, mapWidth, mapHeight)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	gameMap, err := gamemap.ReflectGameMap(gameMiniMap)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	err = network.SendClientReady(serverConn)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	frameTimeout, players, err := network.RecvInitMessage(serverConn)
	fmt.Println(players)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	netclock = network.NewNetclock(int64(frameTimeout), 0.25)

	return StartGame(serverConn, frameTimeout, gameMap, clientName, players)
}

func StartGame(serverConn net.Conn, frameTimeout int, gameMap gamemap.GameMap, clientName string, recievedPlayers []models.Player) int {
	players := make([]*models.Player, len(recievedPlayers))
	for i := 0; i < len(recievedPlayers); i++ {
		players[i] = &recievedPlayers[i]
	}
	ind, isIn := models.NameInPlayers(clientName, recievedPlayers)
	if !isIn {
		fmt.Printf("Can't find player with name %s in recieved players!\n", clientName)
		return 1
	}

	controlledPlayer := players[ind]
	controlledPlayer.Controlled = true

	term, err := terminal.NewTerminal(&gameMap, players)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	defer term.Close()

	var event termbox.Event

	go handleKeypress(serverConn, controlledPlayer, &event, term)
	go handleServerKeys(serverConn, players)

	for gameMap.FoodCount > 0 {
		term.Render()
		if event.Ch == 'q' || event.Key == termbox.KeyCtrlC {
			return 0
		}
		logic.ProcessFrame(players, &gameMap)

		time.Sleep(time.Until(netclock.NextFrameEnd))
	}

	term.RenderWin()
	time.Sleep(time.Second * 5)
	return 0
}

func handleKeypress(serverConn net.Conn, player *models.Player, e *termbox.Event, term *terminal.Terminal) {
	var err error
	for {
		*e = term.Poll()
		if e.Ch == 0 {
			switch e.Key {
			case termbox.KeyArrowUp:
				netclock.WaitUntilSafeFrame()
				err = network.SendClientKey(serverConn, models.MoveUp)
				player.SetDirection(models.MoveUp)
			case termbox.KeyArrowDown:
				netclock.WaitUntilSafeFrame()
				err = network.SendClientKey(serverConn, models.MoveDown)
				player.SetDirection(models.MoveDown)
			case termbox.KeyArrowLeft:
				netclock.WaitUntilSafeFrame()
				err = network.SendClientKey(serverConn, models.MoveLeft)
				player.SetDirection(models.MoveLeft)
			case termbox.KeyArrowRight:
				netclock.WaitUntilSafeFrame()
				err = network.SendClientKey(serverConn, models.MoveRight)
				player.SetDirection(models.MoveRight)
			}
		}
		if err != nil {
			fmt.Println("Some error with recieving keys")
			os.Exit(0)
		}
	}
}

func handleServerKeys(serverConn net.Conn, players []*models.Player) {
	for {
		dir, name, err := network.RecvServerKey(serverConn)
		if err != nil {
			fmt.Println("Error in RecvServerKey function", err)
			return
		}
		id, isIn := models.NameInLPlayers(name, players)
		if !isIn {
			fmt.Printf("Error in RecvServerKey function. Can't find %s in players.\n", name)
			return
		}
		players[id].SetDirection(dir)
	}
}
