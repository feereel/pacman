package terminal

import (
	"github.com/feereel/pacman/internal/gamemap"
	"github.com/nsf/termbox-go"
)

type VisualSymbol struct {
	Fg  termbox.Attribute
	Bg  termbox.Attribute
	Val rune
}

type TerminalElement int

const (
	ControlledPlayer   TerminalElement = iota
	UncontrolledPlayer TerminalElement = iota
	OutBoundWall       TerminalElement = iota
	InBoundWall        TerminalElement = iota
	Food               TerminalElement = iota
	Empty              TerminalElement = iota
	TextTitle          TerminalElement = iota
	TextContent        TerminalElement = iota
	TextPlayer         TerminalElement = iota
)

var GameMapElement = map[gamemap.Cell]TerminalElement{
	gamemap.Empty:  Empty,
	gamemap.Player: UncontrolledPlayer,
	gamemap.Food:   Food,
	gamemap.Wall:   InBoundWall,
}

var ElementToVisual = map[TerminalElement]VisualSymbol{
	ControlledPlayer:   {Fg: termbox.ColorLightMagenta, Bg: termbox.ColorDefault, Val: '@'},
	UncontrolledPlayer: {Fg: termbox.ColorWhite, Bg: termbox.ColorDefault, Val: '@'},
	OutBoundWall:       {Fg: termbox.ColorRed, Bg: termbox.ColorDefault, Val: '#'},
	InBoundWall:        {Fg: termbox.ColorWhite, Bg: termbox.ColorDefault, Val: '#'},
	Food:               {Fg: termbox.ColorLightYellow, Bg: termbox.ColorDefault, Val: '.'},
	Empty:              {Fg: termbox.ColorLightYellow, Bg: termbox.ColorDefault, Val: ' '},
	TextTitle:          {Fg: termbox.ColorBlue, Bg: termbox.ColorDefault},
	TextContent:        {Fg: termbox.ColorWhite, Bg: termbox.ColorDefault},
	TextPlayer:         {Fg: termbox.ColorGreen, Bg: termbox.ColorDefault},
}
