package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/feereel/pacman/cmd/client"
	"github.com/feereel/pacman/cmd/server"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [file]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var showHelp, runAsServer, onlyServerMode bool
	var port, playersCount int
	var ip, playerName string
	mapWidth, mapHeight, mapOccupancy := 20, 15, 0.3

	flag.BoolVar(&showHelp, "h", false, "print usage")
	flag.BoolVar(&runAsServer, "s", false, "run in server mode")
	flag.BoolVar(&onlyServerMode, "os", false, "if set to server mode, the server will not connect to itself")
	flag.IntVar(&port, "p", 4444, "port of a server")
	flag.IntVar(&playersCount, "n", 2, "players count")
	flag.StringVar(&ip, "ip", "127.0.0.1", "ip of a server")
	flag.StringVar(&playerName, "name", "anonymous", "name of created player")
	flag.Usage = usage
	flag.Parse()
	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if runAsServer {
		os.Exit(server.Run(port, playersCount, mapWidth, mapHeight, float32(mapOccupancy), playerName, onlyServerMode))
	} else {
		os.Exit(client.Run(playerName, ip, port, mapWidth, mapHeight))
	}
}
