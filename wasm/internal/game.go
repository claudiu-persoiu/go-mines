package internal

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/claudiu-persoiu/go-mines/internal/elements"
	"github.com/claudiu-persoiu/go-mines/internal/renderer"
)

func ResetGame(level *Level, markMode bool) *Game {
	return NewGame(level, markMode)
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
	GamePaused
)

type Game struct {
	status          GameStatus
	Level           *Level
	marked          int
	menu            *Menu
	eventsHandler   *EventsHandler
	elementsHandler *elements.Handler
	events          chan event
	markMode        bool
	ticker          *time.Ticker
	time            int
	renderer        *renderer.Html
}

func NewGame(level *Level, markMode bool) *Game {
	events := make(chan event)

	g := &Game{
		status:          GameNew,
		Level:           level,
		marked:          0,
		menu:            NewMenu(),
		eventsHandler:   NewEventsHandler(events),
		elementsHandler: elements.NewElementsHandler(level.X, level.Y),
		events:          events,
		time:            0,
		renderer:        renderer.NewHtml(),
		markMode:        markMode,
	}
	g.menu.HideMenu()
	g.GenerateCanvas()
	g.processEvents()
	g.displayTime()

	return g
}

func (g *Game) GenerateCanvas() {
	g.renderer.GenerateCanvas(
		g.status == GameActive || g.status == GameNew,
		g.elementsHandler,
		g.Level.X, g.Level.Y,
		js.FuncOf(g.eventsHandler.EventDown), js.FuncOf(g.eventsHandler.EventUp),
	)

	g.UpdateBombLabel()
}

func (g *Game) UpdateBombLabel() {
	g.renderer.UpdateBombLabel(g.marked, g.Level.Bombs)
}

func (g *Game) ToggleMarkMode() bool {
	g.markMode = !g.markMode
	return g.markMode
}

func (g *Game) processEvents() {
	go func() {
		for e := range g.events {
			switch e.action {
			case "left":
				if g.markMode && g.status != GameOver && g.status != GameNew {
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
			}

			if e.action == "highlight" {
				g.elementsHandler.Highlight(e.key)
				fmt.Println("Highlight on", e.key)
			} else {
				g.elementsHandler.ClearHighlight(e.key)
				fmt.Println("clear highlights", e.key)
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
		g.elementsHandler.GenerateElements(key, g.Level.Bombs)
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
	t := time.NewTicker(time.Second)
	go func() {
		for range t.C {
			g.time++
			g.displayTime()
		}
	}()
	g.ticker = t
}

func (g *Game) displayTime() {
	g.renderer.DisplayTime(g.time)
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

func (g *Game) Pause() {
	if g.status == GamePaused {
		g.status = GameActive
		g.initInterval()
		g.menu.PauseOff()
		g.renderer.UnpauseGame()
	} else if g.status == GameActive {
		g.status = GamePaused
		g.ticker.Stop()
		g.menu.PauseOn()
		g.renderer.PauseGame()
	}
}
