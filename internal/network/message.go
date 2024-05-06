package network

type Serializable interface {
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
}

type Endian int

const (
	BigEndian    Endian = iota
	LittleEndian Endian = iota
)

var NetworkEndian Endian = BigEndian

const (
	PackageClientInitial  uint32 = 0x01
	PackageServerMap      uint32 = 0x10
	PackageClientReady    uint32 = 0x02
	PackageGameStart      uint32 = 0x20
	PackageClientKeyboard uint32 = 0x00
	PackageServerKeyboard uint32 = 0xffffffff
)
