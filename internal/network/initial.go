package network

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/feereel/pacman/internal/models"
)

func SerializeInitMessage(frameTimeout int, players []models.Player) ([]byte, error) {
	if len(players) < 1 {
		return nil, errors.New("0 players were specified")
	}

	data := make([]byte, 8)
	if NetworkEndian == BigEndian {
		binary.BigEndian.PutUint32(data[0:4], uint32(frameTimeout))
		binary.BigEndian.PutUint32(data[4:8], uint32(len(players)))
	} else {
		binary.LittleEndian.PutUint32(data[0:4], uint32(frameTimeout))
		binary.LittleEndian.PutUint32(data[4:8], uint32(len(players)))
	}
	for _, p := range players {
		playerData, err := SerializePlayer(p)
		if err != nil {
			return nil, err
		}

		data = append(data, playerData...)
	}
	return data, nil
}

func DeserializeInitMessage(conn net.Conn) (int, []models.Player, error) {
	data := make([]byte, 8)
	n, err := conn.Read(data)
	if n != 8 || err != nil {
		return 0, nil, err
	}

	var frameTimeout, playersCount int
	if NetworkEndian == BigEndian {
		frameTimeout = int(binary.BigEndian.Uint32(data[0:4]))
		playersCount = int(binary.BigEndian.Uint32(data[4:8]))
	} else {
		frameTimeout = int(binary.LittleEndian.Uint32(data[0:4]))
		playersCount = int(binary.LittleEndian.Uint32(data[4:8]))
	}

	players := make([]models.Player, playersCount)
	for i := 0; i < playersCount; i++ {
		players[i], err = DeserializePlayer(conn)
		if err != nil {
			return frameTimeout, players, err
		}
	}
	return frameTimeout, players, nil
}
