package internal

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"
)

func ResetGame(level *Level) *Game {
	return NewGame(level)
}

type event struct {
	key    string
	action string
}

type GameStatus int

const (
	GameNew GameStatus = iota
	GameActive
	GameOver
	GameReset
)

type Game struct {
	status          GameStatus
	Level           *Level
	marked          int
	menu            *Menu
	canvas          js.Value
	statusElement   js.Value
	counterElement  js.Value
	eventsHandler   *EventsHandler
	elementsHandler *ElementsHandler
	events          chan event
	markMode        bool
	ticker          *time.Ticker
}

func NewGame(level *Level) *Game {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "canvas")
	statusElement := document.Call("getElementById", "status")
	counter := document.Call("getElementById", "counter")
	events := make(chan event)

	g := &Game{
		status:          GameNew,
		Level:           level,
		marked:          0,
		menu:            NewMenu(),
		canvas:          canvas,
		statusElement:   statusElement,
		counterElement:  counter,
		eventsHandler:   NewEventsHandler(events),
		elementsHandler: NewElementsHandler(level),
		events:          events,
	}
	g.menu.HideMenu()
	g.GenerateCanvas()
	g.processEvents()

	return g
}

func falseFunction(this js.Value, args []js.Value) interface{} {
	return false
}

// TODO move html generation to a different structure
func (g *Game) GenerateCanvas() {
	g.canvas.Set("innerHTML", "")
	if g.status == GameOver || g.status == GameReset {
		g.canvas.Set("className", "translucent")
	} else {
		g.canvas.Set("className", "")
	}

	document := js.Global().Get("document")

	df := document.Call("createDocumentFragment")
	table := document.Call("createElement", "table")
	table.Set("id", "table-elements")
	df.Call("appendChild", table)

	tbody := document.Call("createElement", "tbody")
	table.Call("appendChild", tbody)

	for i := 0; i < g.Level.X; i++ {
		tr := document.Call("createElement", "tr")
		tr.Set("oncontextmenu", js.FuncOf(falseFunction))
		tr.Set("onclick", js.FuncOf(falseFunction))
		tbody.Call("appendChild", tr)
		for j := 0; j < g.Level.Y; j++ {
			status := g.elementsHandler.GetElementStatus(arrayToKey(i, j))
			nb := g.elementsHandler.GetNeighbours(arrayToKey(i, j))

			td := document.Call("createElement", "td")

			if status == "empty" && nb > 0 {
				td.Get("style").Set("color", getTextColor(nb))
				td.Set("innerHTML", nb)
			} else {
				td.Set("innerHTML", "&nbsp;")
			}
			if g.status == GameOver && g.elementsHandler.IsBomb(arrayToKey(i, j)) {
				td.Set("className", "exploded")
			} else {
				td.Set("className", g.elementsHandler.GetElementStatus(arrayToKey(i, j)))
			}
			td.Set("id", arrayToKey(i, j))
			td.Call("setAttribute", "unselectable", "on")

			td.Set("onclick", js.FuncOf(falseFunction))
			if g.status != GameOver && g.status != GameReset {
				td.Set("onmousedown", js.FuncOf(g.eventsHandler.EventDown))
				td.Set("onmouseup", js.FuncOf(g.eventsHandler.EventUp))
			} else {
				td.Set("onmousedown", js.FuncOf(falseFunction))
				td.Set("onmouseup", js.FuncOf(falseFunction))
			}
			td.Set("ondblclick", js.FuncOf(falseFunction))
			td.Set("oncontextmenu", js.FuncOf(falseFunction))

			tr.Call("appendChild", td)
		}
	}

	g.canvas.Call("appendChild", df)

	g.UpdateBombLabel()
}

func (g *Game) UpdateBombLabel() {
	g.statusElement.Set("innerHTML", strconv.Itoa(g.marked)+"/"+strconv.Itoa(g.Level.Bombs))
}

func (g *Game) ToggleMarkMode() bool {
	g.markMode = !g.markMode
	return g.markMode
}

func arrayToKey(x, y int) string {
	return strconv.Itoa(x) + "x" + strconv.Itoa(y)
}

func keyToArray(key string) (int, int) {
	var x, y int
	_, err := fmt.Sscanf(key, "%dx%d", &x, &y)
	if err != nil {
		return 0, 0
	}
	return x, y
}

func (g *Game) processEvents() {
	go func() {
		for e := range g.events {
			switch e.action {
			case "left":
				if g.markMode && g.status != GameOver {
					g.markBomb(e.key)
				} else {
					g.revealElement(e.key)
				}
				g.showMarked(e.key)
				fmt.Println("Left click on", e.key)
			case "right":
				g.markBomb(e.key)
				fmt.Println("Right click on", e.key)
			case "both":
				g.showMarked(e.key)
				fmt.Println("Both click on", e.key)
			case "highlight":
				fmt.Println("Highlight on", e.key)
			}
			g.checkFinished()
			g.GenerateCanvas()
		}
	}()
}

func (g *Game) markBomb(key string) {
	mb := g.elementsHandler.MarkBomb(key)
	if mb == "marked" {
		g.marked++
	} else if mb == "new" {
		g.marked--
	}
}

func (g *Game) revealElement(key string) {
	if g.status == GameNew {
		g.status = GameActive
		g.elementsHandler.generateElements(keyToArray(key))
		g.initInterval()
	}

	if g.elementsHandler.IsBomb(key) {
		g.gameOver()
		return
	}

	g.elementsHandler.SetStatus(key, "empty")

	fmt.Println(key, "neighbours:", g.elementsHandler.GetNeighbours(key))
	if g.elementsHandler.GetNeighbours(key) == 0 {
		g.elementsHandler.ClearNeighbourElements(key)
	}

	g.GenerateCanvas()
}

func (g *Game) showMarked(key string) {
	sm := g.elementsHandler.ShowMarked(key)
	if !sm {
		g.gameOver()
	}
}

func (g *Game) gameOver() {
	g.ticker.Stop()
	g.status = GameOver

	g.GenerateCanvas()
	g.menu.ShowMenu("You died ... :(", "reset")
}

func (g *Game) initInterval() {
	s := 0
	t := time.NewTicker(time.Second)
	go func() {
		for range t.C {
			s++
			l := strconv.Itoa(s/60) + ":"
			sec := s % 60
			if sec < 10 {
				l += "0"
			}
			l += strconv.Itoa(sec)
			g.counterElement.Set("innerHTML", l)
		}
	}()
	g.ticker = t
}

func (g *Game) checkFinished() {

	if !g.elementsHandler.CheckFinished() {
		return
	}
	g.ticker.Stop()
	g.status = GameOver

	g.menu.ShowMenu("You win! :)", "reset")
}

func (g *Game) Reset() {
	g.ticker.Stop()
	g.status = GameReset
	g.GenerateCanvas()
	g.menu.ShowMenu("Start fresh?", "reset")
}

func getTextColor(n int) string {
	switch n {
	case 1:
		return "#000000"
	case 2:
		return "#0000FF"
	case 3:
		return "#00FFFF"
	case 4:
		return "#00FF00"
	case 5:
		return "#00FF00"
	case 6:
		return "#00FF00"
	case 7:
		return "#FF0000"
	}
	return "#000000"
}
