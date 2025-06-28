package series

// stringElements is the concrete implementation of Elements for String elements.
type stringElements []stringElement

func (e stringElements) Len() int           { return len(e) }
func (e stringElements) Elem(i int) Element { return &e[i] }

type stringElement struct {
	val string
}

func (e stringElement) Type() Type {
	return String
}

func (e stringElement) Val() any {
	return e.val
}

func (e *stringElement) Set(value any) {

	switch val := value.(type) {
	case string:
		e.val = string(val)
	default:
		return
	}
}

func (e stringElement) ToString() string {
	return e.val
}
