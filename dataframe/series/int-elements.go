package series

import "fmt"

// intElements is the concrete implementation of Elements for Int elements.
type intElements []intElement

func (e intElements) Len() int           { return len(e) }
func (e intElements) Elem(i int) Element { return &e[i] }

type intElement struct {
	val int
}

func (e intElement) Type() Type {
	return String
}

func (e intElement) Val() any {
	return e.val
}

func (e *intElement) Set(value any) {

	switch val := value.(type) {
	case int:
		e.val = int(val)
	default:
		return
	}
}

func (e intElement) ToString() string {
	return fmt.Sprintf("%v", e.val)
}
