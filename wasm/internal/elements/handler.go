package elements

import (
	"fmt"
	"math/rand"
	"slices"
	"strconv"
)

type Handler struct {
	x        int
	y        int
	elements map[string]*Element
}

func NewElementsHandler(x, y int) *Handler {
	eh := &Handler{
		elements: make(map[string]*Element),
		x:        x,
		y:        y,
	}
	eh.resetElements()
	return eh
}

func (eh *Handler) MarkBomb(key string) string {
	if eh.elements[key].GetStatus() == "marked" {
		eh.elements[key].SetStatus("new")
	} else if eh.elements[key].GetStatus() == "new" {
		eh.elements[key].SetStatus("marked")
	}

	return eh.elements[key].GetStatus()
}

func (eh *Handler) IsBomb(key string) bool {
	return eh.elements[key].IsBomb()
}

func (eh *Handler) GetElementStatus(key string) string {
	return eh.elements[key].GetStatus()
}

func (eh *Handler) SetStatus(key, status string) {
	eh.elements[key].SetStatus(status)
}

func (eh *Handler) GetNeighbours(key string) int {
	return eh.elements[key].neighbors
}

func (eh *Handler) ClearNeighbourElements(key string) {
	x, y := keyToArray(key)
	eh.clearNeighbors(x, y, make(map[string]bool))
}

func (eh *Handler) resetElements() {
	eh.elements = make(map[string]*Element)

	for x := 0; x < eh.x; x++ {
		for y := 0; y < eh.y; y++ {
			eh.elements[ArrayToKey(x, y)] = NewElement()
		}
	}
}

func (eh *Handler) GenerateElements(key string, bombs int) {
	x, y := keyToArray(key)
	excludePositions := getNeighborKeys(x, y)
	excludePositions = append(excludePositions, key)

	var keys []string

	for x := 0; x < eh.x; x++ {
		for y := 0; y < eh.y; y++ {
			keys = append(keys, ArrayToKey(x, y))
		}
	}

	for i := 0; i < bombs; i++ {
		r := rand.Intn(len(keys))
		if slices.Contains(excludePositions, keys[r]) {
			i--
			continue
		}
		eh.elements[keys[r]].SetBomb(true)
		keys = slices.Delete(keys, r, r+1)
	}

	for key, element := range eh.elements {
		element.SetNeighbors(eh.getNeighborBombsCount(key))
	}

}

func (eh *Handler) getNeighborBombsCount(key string) int {
	x, y := keyToArray(key)
	count := 0

	neighborKeys := getNeighborKeys(x, y)

	for _, nKey := range neighborKeys {
		if _, exists := eh.elements[nKey]; !exists {
			continue
		}
		if eh.elements[nKey].IsBomb() {
			count++
		}
	}
	return count
}

func (eh *Handler) CheckFinished() bool {
	for _, element := range eh.elements {
		if element.IsBomb() == false && (element.GetStatus() == "marked" || element.GetStatus() == "new") {
			return false
		}
	}
	return true
}

func getNeighborKeys(x, y int) []string {
	return []string{
		ArrayToKey(x-1, y-1),
		ArrayToKey(x-1, y),
		ArrayToKey(x-1, y+1),
		ArrayToKey(x, y-1),
		ArrayToKey(x, y+1),
		ArrayToKey(x+1, y-1),
		ArrayToKey(x+1, y),
		ArrayToKey(x+1, y+1),
	}
}

func (eh *Handler) clearNeighbors(x, y int, emptyNeighbors map[string]bool) {

	if _, exists := emptyNeighbors[ArrayToKey(x, y)]; exists {
		return
	}

	emptyNeighbors[ArrayToKey(x, y)] = true

	keys := getNeighborKeys(x, y)
	for _, key := range keys {
		if _, ok := eh.elements[key]; !ok {
			continue
		}
		if eh.elements[key].IsBomb() || eh.elements[key].GetStatus() == "marked" {
			continue
		}
		eh.elements[key].SetStatus("empty")

		if eh.elements[key].GetNeighbors() == 0 {
			x1, y1 := keyToArray(key)
			eh.clearNeighbors(x1, y1, emptyNeighbors)
		}
	}

	return
}

func (eh *Handler) ShowMarked(key string) bool {
	if eh.GetElementStatus(key) != "empty" || eh.GetNeighbours(key) == 0 {
		return true
	}

	nb := eh.GetNeighbours(key)

	x, y := keyToArray(key)
	for _, neighborKey := range getNeighborKeys(x, y) {
		if _, ok := eh.elements[neighborKey]; !ok {
			continue
		}

		if eh.GetElementStatus(neighborKey) == "marked" {
			if eh.IsBomb(neighborKey) {
				nb--
				continue
			}
			return false
		}
	}

	if nb == 0 {
		eh.ClearNeighbourElements(key)
	}

	return true
}

func (eh *Handler) Highlight(key string) {
	x, y := keyToArray(key)
	for _, neighborKey := range getNeighborKeys(x, y) {
		if _, ok := eh.elements[neighborKey]; !ok {
			continue
		}
		if eh.elements[neighborKey].GetStatus() == "new" {
			eh.elements[neighborKey].SetMarked(true)
		}
	}
}
func (eh *Handler) ClearHighlight(key string) {
	x, y := keyToArray(key)
	for _, neighborKey := range getNeighborKeys(x, y) {
		if _, ok := eh.elements[neighborKey]; !ok {
			continue
		}

		eh.elements[neighborKey].SetMarked(false)
	}
}

func (eh *Handler) IsMarked(key string) bool {
	return eh.elements[key].IsMarked()
}

func ArrayToKey(x, y int) string {
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
