package series

import "fmt"

type floatElement struct {
	val float64
}

func newFloatElement(val float64) *floatElement {
	return &floatElement{val: val}
}

func (e floatElement) Type() Type {
	return String
}

func (e floatElement) Val() any {
	return e.val
}

func (e *floatElement) Set(value any) {

	switch val := value.(type) {
	case float64:
		e.val = float64(val)
	default:
		return
	}
}

func (e floatElement) ToString() string {
	return fmt.Sprintf("%v", e.val)
}
