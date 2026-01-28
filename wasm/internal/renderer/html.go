package renderer

import (
	"strconv"
	"syscall/js"

	"github.com/claudiu-persoiu/go-mines/internal/elements"
)

type Html struct {
	canvas         js.Value
	statusElement  js.Value
	counterElement js.Value
}

func NewHtml() *Html {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "canvas")
	statusElement := document.Call("getElementById", "status")
	counter := document.Call("getElementById", "counter")

	return &Html{
		canvas:         canvas,
		statusElement:  statusElement,
		counterElement: counter,
	}
}

func (r *Html) DisplayTime(time int) {
	l := strconv.Itoa(time/60) + ":"
	sec := time % 60
	if sec < 10 {
		l += "0"
	}
	l += strconv.Itoa(sec)
	r.counterElement.Set("innerHTML", l)
}

func (r *Html) UpdateBombLabel(marked, total int) {
	r.statusElement.Set("innerHTML", strconv.Itoa(marked)+"/"+strconv.Itoa(total))
}

func (r *Html) GenerateCanvas(
	inGame bool,
	elementsHandler *elements.Handler,
	x, y int,
	eventDown js.Func, eventUp js.Func) {

	if !inGame {
		r.canvas.Set("className", "translucent")
	} else {
		r.canvas.Set("className", "")
	}

	document := js.Global().Get("document")

	df := document.Call("createDocumentFragment")
	table := document.Call("createElement", "table")
	table.Set("id", "table-elements")
	df.Call("appendChild", table)

	tbody := document.Call("createElement", "tbody")
	table.Call("appendChild", tbody)

	for i := 0; i < x; i++ {
		tr := document.Call("createElement", "tr")
		tr.Set("oncontextmenu", js.FuncOf(falseFunction))
		tr.Set("onclick", js.FuncOf(falseFunction))
		tbody.Call("appendChild", tr)
		for j := 0; j < y; j++ {
			key := elements.ArrayToKey(i, j)
			status := elementsHandler.GetElementStatus(key)
			nb := elementsHandler.GetNeighbours(key)

			td := document.Call("createElement", "td")

			if status == "empty" && nb > 0 {
				td.Get("style").Set("color", getTextColor(nb))
				td.Set("innerHTML", nb)
			} else {
				td.Set("innerHTML", "&nbsp;")
			}
			if !inGame && elementsHandler.IsBomb(key) {
				td.Set("className", "exploded")
			} else {
				td.Set("className", elementsHandler.GetElementStatus(key))
			}
			td.Set("id", key)
			td.Call("setAttribute", "unselectable", "on")

			td.Set("onclick", js.FuncOf(falseFunction))
			if inGame {
				td.Set("onmousedown", eventDown)
				td.Set("onmouseup", eventUp)
			} else {
				td.Set("onmousedown", js.FuncOf(falseFunction))
				td.Set("onmouseup", js.FuncOf(falseFunction))
			}

			if elementsHandler.IsMarked(key) {
				td.Get("style").Set("borderColor", "#9D9392")
			}

			td.Set("ondblclick", js.FuncOf(falseFunction))
			td.Set("oncontextmenu", js.FuncOf(falseFunction))

			tr.Call("appendChild", td)
		}
	}

	r.canvas.Set("innerHTML", "")
	r.canvas.Call("appendChild", df)
}

func falseFunction(this js.Value, args []js.Value) interface{} {
	return false
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

func (r *Html) PauseGame() {
	r.canvas.Get("style").Set("visibility", "hidden")
}

func (r *Html) UnpauseGame() {
	r.canvas.Get("style").Set("visibility", "visible")
}
