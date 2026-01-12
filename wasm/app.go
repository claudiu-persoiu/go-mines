package main

import (
	"syscall/js"

	"github.com/claudiu-persoiu/go-mines/internal"
)

func main() {

	var _ *internal.Game

	c := make(chan struct{})

	window := js.Global().Get("window")
	window.Set("onload", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		level := internal.GetLevelValues()

		if level.X > 0 && level.Y > 0 && level.Bombs > 0 {
			_ = internal.ResetGame(level)
		}

		return nil
	}))

	<-c
}
