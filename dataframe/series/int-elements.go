package series

import "fmt"

type intElement struct {
	val int
}

func newIntElement(val int) *intElement {
	return &intElement{val: val}
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
