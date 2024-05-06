package server

import (
	"fmt"
	"net"

	"github.com/feereel/pacman/internal/gamemap"
	"github.com/feereel/pacman/internal/models"
	"github.com/feereel/pacman/internal/network"
	"github.com/feereel/pacman/internal/utility"
)

func Run(port, playersCount, mapWidth, mapHeight int, mapOccupancy float32) int {
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

	fmt.Println("Server is listening...")
	for len(players) < playersCount {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}

		fmt.Println("\tA new client is trying to connect...")

		player, err := handleHandshake(conn, gameMap)
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

		if len(players) == 1 {
			player.Position.X = mapWidth - 1
			player.Position.Y = mapHeight - 1
		}

		players = append(players, player)
		fmt.Printf("\tA new client connected with name: %s\n", player.Name)
	}

	frameTimeout := 100

	for _, p := range players {
		err := network.SendInitMessage(p.Conn, frameTimeout, players)
		if err != nil {
			fmt.Printf("There was an error trying to send the initial message. Maybe someone disconnected. %v\n", err)
			return 1
		}
	}

	dirChan := make(chan int, 1)

	for i := 0; i < len(players); i++ {
		go handleKeypress(&players[i], i, dirChan)
	}

	for {
		id := <-dirChan
		player := players[id]
		network.SendServerKey(models.VecToDir[player.Direction], players, player.Name)
		fmt.Printf("Player %s changed direction. Current %v.\n", player.Name, player.Direction)
	}
}

func handleHandshake(conn net.Conn, gameMap gamemap.GameMap) (models.Player, error) {
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

		Position:  utility.Vector2D[int]{X: 0, Y: 0},
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
