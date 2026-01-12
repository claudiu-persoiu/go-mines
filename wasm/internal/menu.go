package internal

import (
	"syscall/js"
)

type Menu struct {
	MenuContainer    js.Value
	MessageContainer js.Value
	PauseScreen      js.Value
	PauseImage       js.Value
	Options          map[string]js.Value
}

func NewMenu() *Menu {
	document := js.Global().Get("document")

	menuContainer := document.Call("getElementById", "menu-screen")
	messageContainer := document.Call("getElementById", "menu-message")
	pauseScreen := document.Call("getElementById", "pause-screen")
	pauseImage := document.Call("getElementById", "pause-image")

	options := make(map[string]js.Value)
	options["reset"] = document.Call("getElementById", "reset-options")
	options["type"] = document.Call("getElementById", "type-options")
	options["custom"] = document.Call("getElementById", "custom-options")

	return &Menu{
		MenuContainer:    menuContainer,
		MessageContainer: messageContainer,
		PauseScreen:      pauseScreen,
		PauseImage:       pauseImage,
		Options:          options,
	}
}

func (m *Menu) HideMenu() {
	m.MenuContainer.Get("style").Set("display", "none")
	for _, option := range m.Options {
		option.Get("style").Set("display", "none")
	}
}

func (m *Menu) ShowMenu(message, menu string) {
	m.HideMenu()
	m.MenuContainer.Get("style").Set("display", "block")
	m.MessageContainer.Set("innerHTML", message)

	if option, exists := m.Options[menu]; exists {
		option.Get("style").Set("display", "block")
	}
	m.PauseScreen.Get("style").Set("display", "none")
}

func (m *Menu) PauseOn() {
	m.PauseScreen.Get("style").Set("display", "block")
	m.PauseImage.Set("className", "play")
}

func (m *Menu) PauseOff() {
	m.PauseScreen.Get("style").Set("display", "none")
	m.PauseImage.Set("className", "pause")
}
