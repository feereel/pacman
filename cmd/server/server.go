package server

import (
	"fmt"
	"net"
	"time"

	"github.com/feereel/pacman/cmd/client"
	"github.com/feereel/pacman/internal/gamemap"
	"github.com/feereel/pacman/internal/models"
	"github.com/feereel/pacman/internal/network"
	"github.com/feereel/pacman/internal/utility"
)

var netclock *network.Netclock

func Run(port, playersCount, mapWidth, mapHeight int, mapOccupancy float32, serverName string, onlyServerMode bool, frameTimeout int) int {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))

	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer listener.Close()

	players := make([]models.Player, 0)

	gameMap, err := gamemap.NewGameMap(mapWidth, mapHeight, mapOccupancy)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	var playerIDToPos = map[int]utility.Vector2D[int]{
		0: {X: 0, Y: 0},
		1: {X: mapWidth*2 - 1, Y: 0},
		2: {X: 0, Y: mapHeight*2 - 1},
		3: {X: mapWidth*2 - 1, Y: mapHeight*2 - 1},
	}

	if !onlyServerMode {
		go client.Run(serverName, "127.0.0.1", port, mapWidth, mapHeight)
	}

	fmt.Println("Server is listening...")
	for len(players) < playersCount {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}

		fmt.Println("\tA new client is trying to connect...")

		player, err := handleHandshake(conn, gameMap, playerIDToPos[len(players)])
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		if _, isIn := models.NameInPlayers(player.Name, players); isIn {
			fmt.Printf("\tName %s reserved.\n", player.Name)
			conn.Close()
			continue
		}

		players = append(players, player)
		fmt.Printf("\tA new client connected with name: %s\n", player.Name)
	}

	for _, p := range players {
		err := network.SendInitMessage(p.Conn, frameTimeout, players)
		if err != nil {
			fmt.Printf("There was an error trying to send the initial message. Maybe someone disconnected. %v\n", err)
			return 1
		}
	}

	netclock = network.NewNetclock(int64(frameTimeout), 0.25)

	dirChan := make(chan int, 1)

	for i := 0; i < len(players); i++ {
		go handleKeypress(&players[i], i, dirChan)
	}

	go handleKeysend(players, dirChan)

	time.Sleep(time.Second * 100000000)
	return 0
}

func handleHandshake(conn net.Conn, gameMap gamemap.GameMap, spawnPos utility.Vector2D[int]) (models.Player, error) {
	name, err := network.RecvPlayerName(conn)
	if err != nil {
		return models.Player{}, err
	}

	err = network.SendGameMap(conn, gameMap)
	if err != nil {
		return models.Player{}, err
	}

	err = network.RecvClientReady(conn)
	if err != nil {
		return models.Player{}, err
	}

	return models.Player{
		Name:       name,
		Score:      0,
		Controlled: false,
		Conn:       conn,

		Position:  spawnPos,
		Direction: models.DirToVec[models.MoveDown],
	}, nil
}

func handleKeypress(player *models.Player, id int, c chan int) {
	for {
		dir, err := network.RecvClientKey(player.Conn)
		if err != nil {
			fmt.Printf("Something happend with %v client (%s). Error: %v. \n", id, player.Name, err)
			return
		}

		c <- id
		player.SetDirection(dir)
	}
}

func handleKeysend(players []models.Player, c chan int) {
	netclock.WaitUntilSafeFrame()
	for {
		player := players[<-c]
		network.SendServerKey(models.VecToDir[player.Direction], players, player.Name)
	}
}
