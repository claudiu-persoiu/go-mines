package elements

type Element struct {
	isBomb    bool
	status    string
	modified  bool
	neighbors int
	marked    bool
}

func NewElement() *Element {
	return &Element{
		isBomb:    false,
		status:    "new",
		modified:  false,
		neighbors: 0,
		marked:    false,
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

func (e *Element) IsMarked() bool {
	return e.marked
}

func (e *Element) SetMarked(marked bool) {
	e.marked = marked
}

func (e *Element) GetNeighbors() int {
	return e.neighbors
}

func (e *Element) SetNeighbors(neighbors int) {
	e.neighbors = neighbors
}
