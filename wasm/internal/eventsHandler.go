package internal

import (
	"syscall/js"
)

type EventsHandler struct {
	leftClick   bool
	rightClick  bool
	middleClick bool
	events      chan event
}

func NewEventsHandler(events chan event) *EventsHandler {
	return &EventsHandler{
		leftClick:   false,
		rightClick:  false,
		middleClick: false,
		events:      events,
	}
}

func (eh *EventsHandler) EventUp(this js.Value, args []js.Value) interface{} {

	var clicked string

	e := args[0]

	if !e.Get("button").IsUndefined() {
		switch e.Get("button").Int() {
		case 0:
			eh.leftClick = false
			clicked = "left"
		case 1:
			eh.leftClick = false
			clicked = "left"
		case 2:
			eh.rightClick = false
			clicked = "right"
		case 3:
			eh.leftClick = false
			eh.rightClick = false
			clicked = "both"
		case 4:
			eh.middleClick = false
			clicked = "both"
		}
	} else if !e.Get("which").IsUndefined() {
		switch e.Get("which").Int() {
		case 1:
			eh.leftClick = false
			clicked = "left"
		case 2:
			eh.middleClick = false
			clicked = "both"
		case 3:
			eh.rightClick = false
			clicked = "right"
		}
	}

	// handle both clicked
	if clicked == "left" && eh.rightClick {
		eh.rightClick = false
		clicked = "both"
	} else if clicked == "right" && eh.leftClick {
		eh.leftClick = false
		clicked = "both"
	}

	key := this.Get("id").String()

	eh.events <- event{
		key:    key,
		action: clicked,
	}

	//
	//switch clicked {
	//case "left":
	//	fmt.Println("Left click on", key)
	//case "right":
	//	fmt.Println("Right click on", key)
	//case "both":
	//	fmt.Println("Both click on", key)
	//}
	//
	//// TODO clear highlights

	return nil
}

func (eh *EventsHandler) EventDown(this js.Value, args []js.Value) interface{} {
	e := args[0]

	if !e.Get("button").IsUndefined() {
		switch e.Get("button").Int() {
		case 0:
			eh.leftClick = true
		case 1:
			eh.leftClick = true
		case 2:
			eh.rightClick = true
		case 3:
			eh.leftClick = true
			eh.rightClick = true
		case 4:
			eh.middleClick = true
		}
	} else if !e.Get("which").IsUndefined() {
		switch e.Get("which").Int() {
		case 1:
			eh.leftClick = true
		case 2:
			eh.middleClick = true
		case 3:
			eh.rightClick = true
		}
	}

	if eh.middleClick || eh.leftClick {
		key := this.Get("id").String()
		eh.events <- event{
			key:    key,
			action: "highlight",
		}
	}

	return nil
}
