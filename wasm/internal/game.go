package internal

import (
	"fmt"
	"strconv"
	"syscall/js"
)

func ResetGame(level *Level) *Game {
	return NewGame(level)
}

type event struct {
	key    string
	action string
}

type Game struct {
	Level           *Level
	marked          int
	menu            *Menu
	canvas          js.Value
	statusElement   js.Value
	eventsHandler   *EventsHandler
	elementsHandler *ElementsHandler
	events          chan event
	markMode        bool
	inGame          bool
}

func NewGame(level *Level) *Game {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "canvas")
	statusElement := document.Call("getElementById", "status")
	events := make(chan event)

	g := &Game{
		Level:           level,
		marked:          0,
		menu:            NewMenu(),
		canvas:          canvas,
		statusElement:   statusElement,
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

func (g *Game) GenerateCanvas() {
	g.canvas.Set("innerHTML", "")
	g.canvas.Set("className", "")

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
			td := document.Call("createElement", "td")
			td.Set("innerHTML", "&nbsp;")
			td.Set("className", g.elementsHandler.GetElementStatus(arrayToKey(i, j)))
			td.Set("id", arrayToKey(i, j))
			td.Call("setAttribute", "unselectable", "on")

			td.Set("onclick", js.FuncOf(falseFunction))
			td.Set("onmousedown", js.FuncOf(g.eventsHandler.EventDown))
			td.Set("onmouseup", js.FuncOf(g.eventsHandler.EventUp))
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
				if g.markMode && g.inGame {
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
			g.GenerateCanvas()
		}
	}()
}

func (g *Game) markBomb(key string) {
	if g.elementsHandler.MarkBomb(key) == "marked" {
		g.marked++
	} else {
		g.marked--
	}
}

func (g *Game) revealElement(key string) {
	// TODO
}

func (g *Game) showMarked(key string) {
	// TODO
}
