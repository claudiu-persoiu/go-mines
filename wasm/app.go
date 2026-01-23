package main

import (
	"strconv"
	"syscall/js"

	"github.com/claudiu-persoiu/go-mines/internal"
)

func main() {

	var _ *internal.Game

	c := make(chan struct{})

	var g *internal.Game

	restartFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		level := internal.GetLevelValues()

		if level.X > 0 && level.Y > 0 && level.Bombs > 0 {
			g = internal.ResetGame(level)
		}

		return nil
	})

	window := js.Global().Get("window")
	window.Set("onload", restartFunc)

	document := js.Global().Get("document")
	document.Call("getElementById", "new-game").Set("onclick", restartFunc)

	sa := document.Call("getElementById", "switch-action")

	sa.Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		m := g.ToggleMarkMode()
		if m {
			sa.Call("getElementsByTagName", "div").Index(0).Set("className", "mark-flag")
		} else {
			sa.Call("getElementsByTagName", "div").Index(0).Set("className", "mark-explode")
		}

		return nil
	}))

	// Change level
	document.Call("getElementById", "different-level").Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		document.Call("getElementById", "reset-options").Get("style").Set("display", "none")
		document.Call("getElementById", "type-options").Get("style").Set("display", "block")

		return nil
	}))

	// Select level
	document.Call("querySelectorAll", ".option-start").Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		option := args[0]
		option.Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			ds := option.Get("dataset")

			x, e1 := strconv.Atoi(ds.Get("x").String())
			y, e2 := strconv.Atoi(ds.Get("y").String())
			bombs, e3 := strconv.Atoi(ds.Get("bombs").String())

			if e1 != nil || e2 != nil || e3 != nil {
				return nil
			}
			level := &internal.Level{X: x, Y: y, Bombs: bombs}
			internal.SetLevelValues(level)
			g = internal.ResetGame(level)

			return nil
		}))

		return nil
	}))

	document.Call("getElementById", "option-custom").Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		document.Call("getElementById", "custom-options").Get("style").Set("display", "block")
		return nil
	}))

	document.Call("getElementById", "custom-options-cancel").Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		document.Call("getElementById", "custom-options").Get("style").Set("display", "none")
		return nil
	}))

	document.Call("getElementById", "option-custom-start").Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		x, e1 := strconv.Atoi(document.Call("getElementById", "custom-x").Get("value").String())
		y, e2 := strconv.Atoi(document.Call("getElementById", "custom-y").Get("value").String())
		bombs, e3 := strconv.Atoi(document.Call("getElementById", "custom-mines").Get("value").String())

		if e1 != nil || e2 != nil || e3 != nil {
			return nil
		}

		level := &internal.Level{X: x, Y: y, Bombs: bombs}
		internal.SetLevelValues(level)
		g = internal.ResetGame(level)

		return nil
	}))

	document.Call("getElementById", "reset").Set("onclick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.Reset()
		return nil
	}))

	// https://egghead.io/lessons/go-call-a-go-webassembly-function-from-javascript

	<-c
}
