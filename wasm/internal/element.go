package internal

type Element struct {
	isBomb    bool
	status    string
	modified  bool
	neighbors int
}

func NewElement() *Element {
	return &Element{
		isBomb:    false,
		status:    "",
		modified:  false,
		neighbors: 0,
	}
}

func (e *Element) IsBomb() bool {
	return e.isBomb
}

func (e *Element) SetBomb(isBomb bool) {
	e.isBomb = isBomb
}

func (e *Element) GetStatus() string {
	return e.status
}

func (e *Element) SetStatus(status string) {
	e.status = status
}

func (e *Element) IsModified() bool {
	return e.modified
}

func (e *Element) SetModified(modified bool) {
	e.modified = modified
}

func (e *Element) GetNeighbors() int {
	return e.neighbors
}

func (e *Element) SetNeighbors(neighbors int) {
	e.neighbors = neighbors
}
