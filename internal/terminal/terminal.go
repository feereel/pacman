package terminal

import (
	"fmt"
	"reflect"

	"github.com/feereel/pacman/internal/gamemap"
	"github.com/feereel/pacman/internal/models"
	"github.com/nsf/termbox-go"
)

type Terminal struct {
	Width, Height    int
	GameMap          *gamemap.GameMap
	Players          []*models.Player
	controlledPlayer *models.Player
}

func NewTerminal(gameMap *gamemap.GameMap, players []*models.Player) (*Terminal, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}

	w, h := termbox.Size()
	var controlledPlayer *models.Player = nil
	for _, p := range players {
		if p.Controlled {
			controlledPlayer = p
			break
		}
	}

	return &Terminal{Width: w, Height: h, GameMap: gameMap, Players: players, controlledPlayer: controlledPlayer}, nil
}

func (t *Terminal) Resize(width, height int) {
	t.Width = width
	t.Height = height
	t.Render()
}

func (t *Terminal) setText(x int, y int, text string, v VisualSymbol) (retX, retY int) {
	runeArray := []rune(text)
	for i := 0; i < len(runeArray); i++ {
		termbox.SetCell(x+i, y, runeArray[i], v.Fg, v.Bg)
	}
	return x + len(runeArray), y
}

func (t *Terminal) renderBorder() {
	for y := 0; y < t.GameMap.Height+2; y++ {
		var v VisualSymbol = ElementToVisual[OutBoundWall]
		termbox.SetCell(0, y, v.Val, v.Fg, v.Bg)
		termbox.SetCell(t.GameMap.Width+1, y, v.Val, v.Fg, v.Bg)
	}
	for x := 0; x < t.GameMap.Width+2; x++ {
		var v VisualSymbol = ElementToVisual[OutBoundWall]
		termbox.SetCell(x, 0, v.Val, v.Fg, v.Bg)
		termbox.SetCell(x, t.GameMap.Height+1, v.Val, v.Fg, v.Bg)
	}
}

func (t *Terminal) renderMap() {
	for y := 0; y < t.GameMap.Height; y++ {
		for x := 0; x < t.GameMap.Width; x++ {
			var v VisualSymbol = ElementToVisual[GameMapElement[t.GameMap.Grid[y][x]]]
			termbox.SetCell(x+1, y+1, v.Val, v.Fg, v.Bg)
		}
	}
}

func (t *Terminal) renderPlayer(x int, y int, player *models.Player, position int) (retX, retY int) {
	var name string
	if len(player.Name) > 15 {
		name = player.Name[:12] + "..."
	} else {
		name = player.Name
	}

	var endX, endY int
	endX, endY = t.setText(x, y, fmt.Sprintf("%v. %-15s", position, name), ElementToVisual[TextPlayer])
	endX, endY = t.setText(endX+2, endY, fmt.Sprintf("%v", player.Score), ElementToVisual[TextContent])
	return endX, endY
}

func (t *Terminal) renderGameInfo() {
	edgeX := t.GameMap.Width + 4

	// Показывает рейтинг пользователей
	_, y := t.setText(edgeX, 0, "Players rating", ElementToVisual[TextTitle])
	for i, player := range t.Players {
		y += 2
		_, y = t.renderPlayer(edgeX, y, player, i)
	}

	// Показывает сколько осталось еды
	x, y := t.setText(edgeX, y+2, "Food left:", ElementToVisual[TextTitle])
	t.setText(x, y, fmt.Sprintf("%v", t.GameMap.FoodCount), ElementToVisual[TextContent])

}

func (t *Terminal) Render() {
	termbox.Clear(termbox.ColorWhite, termbox.ColorDefault)

	t.renderBorder()
	t.renderMap()
	t.sortPlayer()
	if t.Width-t.GameMap.Width > 25 && t.Height > 2*(len(t.Players)+1) {
		t.renderGameInfo()
	}

	if t.controlledPlayer != nil {
		var v VisualSymbol = ElementToVisual[ControlledPlayer]
		termbox.SetCell(t.controlledPlayer.Position.X+1, t.controlledPlayer.Position.Y+1, v.Val, v.Fg, v.Bg)
	}

	termbox.Flush()
}

func (t *Terminal) Poll() termbox.Event {
	for {
		switch e := termbox.PollEvent(); e.Type {
		case termbox.EventKey:
			return e
		case termbox.EventResize:
			t.Resize(e.Width, e.Height)
		}
	}
}

func (t *Terminal) Close() {
	termbox.Close()
}

func (t *Terminal) sortPlayer() {
	swap := reflect.Swapper(t.Players)
	for i := 1; i < len(t.Players); i++ {
		if t.Players[i].Score > t.Players[i-1].Score {
			swap(i-1, i)
		}
	}
}
