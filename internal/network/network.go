package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/feereel/pacman/internal/gamemap"
	"github.com/feereel/pacman/internal/models"
	"github.com/feereel/pacman/internal/utility"
)

func RecvHeader(conn net.Conn, checkType uint32) (int, error) {
	data := make([]byte, 12)
	n, err := conn.Read(data)
	if n != 12 || err != nil {
		return 0, err
	}
	var magic, ptype, datasize uint32
	if NetworkEndian == BigEndian {
		magic = binary.BigEndian.Uint32(data[0:4])
		ptype = binary.BigEndian.Uint32(data[4:8])
		datasize = binary.BigEndian.Uint32(data[8:12])
	} else {
		magic = binary.LittleEndian.Uint32(data[0:4])
		ptype = binary.LittleEndian.Uint32(data[4:8])
		datasize = binary.LittleEndian.Uint32(data[8:12])
	}

	if magic != 0xabcdfe01 {
		s := fmt.Sprintf("wrong magic is sent by user, given: %v, expected %v", magic, 0xabcdfe01)
		return 0, errors.New(s)
	} else if ptype != checkType {
		s := fmt.Sprintf("wrong ptype is sent by user, given: %v, expected %v", checkType, ptype)
		return 0, errors.New(s)
	}
	return int(datasize), nil
}

func SendHeader(conn net.Conn, ptype uint32, datasize uint32) error {
	data := make([]byte, 12)
	if NetworkEndian == BigEndian {
		binary.BigEndian.PutUint32(data[0:4], 0xabcdfe01)
		binary.BigEndian.PutUint32(data[4:8], ptype)
		binary.BigEndian.PutUint32(data[8:12], datasize)
	} else {
		binary.LittleEndian.PutUint32(data[0:4], 0xabcdfe01)
		binary.LittleEndian.PutUint32(data[4:8], ptype)
		binary.LittleEndian.PutUint32(data[8:12], datasize)
	}
	n, err := conn.Write(data)
	if n == 0 || err != nil {
		return err
	}
	return nil
}

func SendPlayerName(conn net.Conn, playerName string) error {
	var data []byte
	if NetworkEndian == BigEndian {
		data = []byte(playerName)
	} else {
		data = utility.ReverseBytes([]byte(playerName))
	}
	if len(data) > 255 {
		return errors.New("name is longer than 255 characters")
	}

	// Send name
	err := SendHeader(conn, PackageClientInitial, uint32(len(data)))
	if err != nil {
		return err
	}

	n, err := conn.Write(data)
	if n == 0 || err != nil {
		return err
	}
	return nil
}

func RecvPlayerName(conn net.Conn) (string, error) {
	size, err := RecvHeader(conn, PackageClientInitial)
	if err != nil {
		return "", err
	}

	// Recv name
	input := make([]byte, size)
	n, err := conn.Read(input)
	if n == 0 || err != nil {
		return "", err
	}
	if NetworkEndian == BigEndian {
		return string(input[0:n]), nil
	} else {
		return string(utility.ReverseBytes(input[0:n])), nil
	}
}

func SendGameMap(conn net.Conn, gameMap gamemap.GameMap) error {
	data, err := SerializeGameMap(gameMap)
	if err != nil {
		return err
	}

	// Send gamemap
	err = SendHeader(conn, PackageServerMap, uint32(len(data)))
	if err != nil {
		return err
	}

	n, err := conn.Write(data)
	if n == 0 || err != nil {
		return err
	}

	return nil
}

func RecvGameMap(conn net.Conn, width int, height int) (gamemap.GameMap, error) {
	_, err := RecvHeader(conn, PackageServerMap)
	if err != nil {
		return gamemap.GameMap{}, err
	}

	// Recv gamemap
	return DeserializeGameMap(conn, width, height)
}

func SendClientReady(conn net.Conn) error {
	// Send ready message
	return SendHeader(conn, PackageClientReady, 0)
}

func RecvClientReady(conn net.Conn) error {
	// Recv ready message
	_, err := RecvHeader(conn, PackageClientReady)
	return err
}

func SendInitMessage(conn net.Conn, frameTimeout int, players []models.Player) error {
	data, err := SerializeInitMessage(frameTimeout, players)
	if err != nil {
		return err
	}

	// Send initial message
	err = SendHeader(conn, PackageGameStart, uint32(len(data)))
	if err != nil {
		return err
	}
	n, err := conn.Write(data)
	if n == 0 || err != nil {
		return err
	}

	return nil
}

func RecvInitMessage(conn net.Conn) (int, []models.Player, error) {
	_, err := RecvHeader(conn, PackageGameStart)
	if err != nil {
		return 0, nil, err
	}

	// Recv initial message
	return DeserializeInitMessage(conn)
}

func SendClientKey(conn net.Conn, direction models.MoveDirection) error {
	// Send client key
	err := SendHeader(conn, PackageClientKeyboard, 1)
	if err != nil {
		return err
	}
	n, err := conn.Write([]byte{byte(direction)})
	if n == 0 || err != nil {
		return err
	}
	return nil
}

func RecvClientKey(conn net.Conn) (models.MoveDirection, error) {
	size, err := RecvHeader(conn, PackageClientKeyboard)
	if err != nil {
		return 0, err
	}

	// Recv client key
	input := make([]byte, size)
	n, err := conn.Read(input)
	if n != 1 || err != nil {
		return 0, err
	}

	return models.MoveDirection(input[0]), nil
}

func SendServerKey(direction models.MoveDirection, players []models.Player, senderName string) error {
	var data []byte
	if NetworkEndian == BigEndian {
		data = []byte(senderName)
	} else {
		data = utility.ReverseBytes([]byte(senderName))

	}

	// Send server key
	for _, p := range players {
		if p.Name == senderName {
			continue
		}
		sendData := append([]byte{byte(direction)}, data...)
		err := SendHeader(p.Conn, PackageServerKeyboard, uint32(len(sendData)))
		if err != nil {
			return err
		}
		n, err := p.Conn.Write(sendData)
		if n == 0 || err != nil {
			return err
		}
	}

	return nil
}

func RecvServerKey(conn net.Conn) (models.MoveDirection, string, error) {
	size, err := RecvHeader(conn, PackageServerKeyboard)
	if err != nil {
		return 0, "", err
	}

	// Recv client key
	input := make([]byte, size)
	n, err := conn.Read(input)
	if n != size || err != nil {
		return 0, "", err
	}

	var name string
	if NetworkEndian == BigEndian {
		name = string(input[1:])
	} else {
		name = string(utility.ReverseBytes(input[1:]))
	}
	return models.MoveDirection(input[0]), name, nil
}
