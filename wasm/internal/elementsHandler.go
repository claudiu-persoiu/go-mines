package internal

import (
	"math/rand"
	"slices"
)

type ElementsHandler struct {
	level    *Level
	elements map[string]*Element
}

func NewElementsHandler(l *Level) *ElementsHandler {
	eh := &ElementsHandler{
		elements: make(map[string]*Element),
		level:    l,
	}
	eh.resetElements()
	return eh
}

func (eh *ElementsHandler) MarkBomb(key string) string {
	if eh.elements[key].GetStatus() == "marked" {
		eh.elements[key].SetStatus("new")
	} else if eh.elements[key].GetStatus() == "new" {
		eh.elements[key].SetStatus("marked")
	}

	return eh.elements[key].GetStatus()
}

func (eh *ElementsHandler) GetElementStatus(key string) string {
	return eh.elements[key].GetStatus()
}

func (eh *ElementsHandler) resetElements() {
	eh.elements = make(map[string]*Element)

	for x := 0; x < eh.level.X; x++ {
		for y := 0; y < eh.level.Y; y++ {
			eh.elements[arrayToKey(x, y)] = NewElement()
		}
	}
}

func (eh *ElementsHandler) generateElements(x, y int) {

	excludePositions := getNeighborKeys(x, y)

	var keys []string

	for x := 0; x < eh.level.X; x++ {
		for y := 0; y < eh.level.Y; y++ {
			keys = append(keys, arrayToKey(x, y))
		}
	}

	for i := 0; i < eh.level.Bombs; i++ {
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

func (eh *ElementsHandler) getNeighborBombsCount(key string) int {
	x, y := keyToArray(key)
	count := 0

	neighborKeys := getNeighborKeys(x, y)

	for _, nKey := range neighborKeys {
		if eh.elements[nKey].IsBomb() {
			count++
		}
	}
	return count
}

func getNeighborKeys(x, y int) []string {
	return []string{
		arrayToKey(x-1, y-1),
		arrayToKey(x-1, y),
		arrayToKey(x-1, y+1),
		arrayToKey(x, y-1),
		arrayToKey(x, y+1),
		arrayToKey(x+1, y-1),
		arrayToKey(x+1, y),
		arrayToKey(x+1, y+1),
	}
}
