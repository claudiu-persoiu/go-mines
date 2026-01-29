package renderer

import (
	"strconv"
	"syscall/js"

	"github.com/claudiu-persoiu/go-mines/internal/level"
)

func GetLevelValues() *level.Level {
	localStorage := js.Global().Get("localStorage")

	x, e1 := strconv.Atoi(localStorage.Call("getItem", "mines-x").String())
	y, e2 := strconv.Atoi(localStorage.Call("getItem", "mines-y").String())
	bombs, e3 := strconv.Atoi(localStorage.Call("getItem", "mines-elements").String())

	if e1 != nil || e2 != nil || e3 != nil {
		return &level.Level{X: 0, Y: 0, Bombs: 0}
	}

	return &level.Level{X: x, Y: y, Bombs: bombs}
}

func SetLevelValues(level *level.Level) {
	localStorage := js.Global().Get("localStorage")

	localStorage.Call("setItem", "mines-x", strconv.Itoa(level.X))
	localStorage.Call("setItem", "mines-y", strconv.Itoa(level.Y))
	localStorage.Call("setItem", "mines-elements", strconv.Itoa(level.Bombs))
}
