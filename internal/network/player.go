package network

import (
	"encoding/binary"
	"net"

	"github.com/feereel/pacman/internal/models"
	"github.com/feereel/pacman/internal/utility"
)

func SerializePlayer(player models.Player) ([]byte, error) {
	name := []byte(player.Name)
	data := make([]byte, 16+len(name))
	if NetworkEndian == BigEndian {
		binary.BigEndian.PutUint32(data[0:4], uint32(player.Position.X))
		binary.BigEndian.PutUint32(data[4:8], uint32(player.Position.Y))
		binary.BigEndian.PutUint32(data[8:12], uint32(models.VecToDir[player.Direction]))
		binary.BigEndian.PutUint32(data[12:16], uint32(len(name)))
		copy(data[16:], name[:])
	} else {
		binary.LittleEndian.PutUint32(data[0:4], uint32(player.Position.X))
		binary.LittleEndian.PutUint32(data[4:8], uint32(player.Position.Y))
		binary.LittleEndian.PutUint32(data[8:12], uint32(models.VecToDir[player.Direction]))
		binary.LittleEndian.PutUint32(data[12:16], uint32(len(name)))
		copy(data[16:], utility.ReverseBytes(name[:]))
	}
	return data, nil
}

func DeserializePlayer(conn net.Conn) (models.Player, error) {
	data := make([]byte, 16)
	n, err := conn.Read(data)
	if n != 16 || err != nil {
		return models.Player{}, err
	}

	var position, direction utility.Vector2D[int]
	var nameLength int
	if NetworkEndian == BigEndian {
		position = utility.Vector2D[int]{X: int(binary.BigEndian.Uint32(data[0:4])), Y: int(binary.BigEndian.Uint32(data[4:8]))}
		direction = models.DirToVec[models.MoveDirection(binary.BigEndian.Uint32(data[8:12]))]
		nameLength = int(binary.BigEndian.Uint32(data[12:16]))
	} else {
		position = utility.Vector2D[int]{X: int(binary.LittleEndian.Uint32(data[0:4])), Y: int(binary.LittleEndian.Uint32(data[4:8]))}
		direction = models.DirToVec[models.MoveDirection(binary.LittleEndian.Uint32(data[8:12]))]
		nameLength = int(binary.LittleEndian.Uint32(data[12:16]))
	}

	data = make([]byte, nameLength)
	n, err = conn.Read(data)
	if n != nameLength || err != nil {
		return models.Player{}, err
	}

	var name string
	if NetworkEndian == BigEndian {
		name = string(data[:])
	} else {
		name = string(utility.ReverseBytes(data[:]))
	}

	return models.Player{
		Name:       name,
		Score:      0,
		Controlled: false,

		Position:  position,
		Direction: direction,
	}, nil
}
