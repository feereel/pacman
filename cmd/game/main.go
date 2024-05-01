package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/feereel/pacman/cmd/client"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [file]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var showHelp bool
	flag.BoolVar(&showHelp, "h", false, "print usage")
	flag.BoolVar(&showHelp, "help", false, "print usage")

	flag.Usage = usage
	flag.Parse()
	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	os.Exit(client.Run())
}
