package main

import (
	"strconv"
	"syscall/js"

	"github.com/claudiu-persoiu/go-mines/internal"
	"github.com/claudiu-persoiu/go-mines/internal/elements"
	"github.com/claudiu-persoiu/go-mines/internal/level"
	"github.com/claudiu-persoiu/go-mines/internal/renderer"
)

func main() {

	var _ *internal.Game

	c := make(chan struct{})

	var g *internal.Game
	markMode := false

	restartGame := func() {
		l := renderer.GetLevelValues()

		if l.X > 0 && l.Y > 0 && l.Bombs > 0 {
			g = internal.ResetGame(l, markMode, elements.NewHandler(l.X, l.Y))
		}
	}

	restartFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		restartGame()
		return nil
	})

	document := js.Global().Get("document")
	document.Call("getElementById", "new-game").Call("addEventListener", "click", restartFunc)

	sa := document.Call("getElementById", "switch-action")

	sa.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		markMode = g.ToggleMarkMode()
		if markMode {
			sa.Call("getElementsByTagName", "div").Index(0).Set("className", "mark-flag")
		} else {
			sa.Call("getElementsByTagName", "div").Index(0).Set("className", "mark-explode")
		}

		return nil
	}))

	// Change level
	document.Call("getElementById", "different-level").Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		document.Call("getElementById", "reset-options").Get("style").Set("display", "none")
		document.Call("getElementById", "type-options").Get("style").Set("display", "block")

		return nil
	}))

	// Select level
	document.Call("querySelectorAll", ".option-start").Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		option := args[0]
		option.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			ds := option.Get("dataset")

			x, e1 := strconv.Atoi(ds.Get("x").String())
			y, e2 := strconv.Atoi(ds.Get("y").String())
			bombs, e3 := strconv.Atoi(ds.Get("bombs").String())

			if e1 != nil || e2 != nil || e3 != nil {
				return nil
			}
			l := &level.Level{X: x, Y: y, Bombs: bombs}
			renderer.SetLevelValues(l)
			g = internal.ResetGame(l, markMode, elements.NewHandler(l.X, l.Y))

			return nil
		}))

		return nil
	}))

	document.Call("getElementById", "option-custom").Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		document.Call("getElementById", "custom-options").Get("style").Set("display", "block")
		return nil
	}))

	document.Call("getElementById", "custom-options-cancel").Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		document.Call("getElementById", "custom-options").Get("style").Set("display", "none")
		return nil
	}))

	document.Call("getElementById", "option-custom-start").Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		x, e1 := strconv.Atoi(document.Call("getElementById", "custom-x").Get("value").String())
		y, e2 := strconv.Atoi(document.Call("getElementById", "custom-y").Get("value").String())
		bombs, e3 := strconv.Atoi(document.Call("getElementById", "custom-mines").Get("value").String())

		if e1 != nil || e2 != nil || e3 != nil {
			return nil
		}

		l := &level.Level{X: x, Y: y, Bombs: bombs}
		renderer.SetLevelValues(l)
		g = internal.ResetGame(l, markMode, elements.NewHandler(l.X, l.Y))

		return nil
	}))

	document.Call("getElementById", "reset").Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.Reset()
		return nil
	}))

	document.Call("querySelectorAll", ".pause-action").Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		option := args[0]
		option.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			g.Pause()
			return nil
		}))
		return nil
	}))

	restartGame()

	<-c
}
