package main

import (
	"syscall/js"

	"github.com/claudiu-persoiu/go-mines/internal"
)

func main() {

	var _ *internal.Game

	c := make(chan struct{})

	restartFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		level := internal.GetLevelValues()

		if level.X > 0 && level.Y > 0 && level.Bombs > 0 {
			_ = internal.ResetGame(level)
		}

		return nil
	})

	window := js.Global().Get("window")
	window.Set("onload", restartFunc)

	document := js.Global().Get("document")
	document.Call("getElementById", "new-game").Set("onclick", restartFunc)

	// https://egghead.io/lessons/go-call-a-go-webassembly-function-from-javascript

	<-c
}
