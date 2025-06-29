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

func (e intElement) ToInt() (int, error) {
	return e.val, nil
}

func (e intElement) ToFloat() float64 {
	return float64(e.val)
}
func (e intElement) ToBool() (bool, error) {
	switch e.val {
	case 1:
		return true, nil
	case 0:
		return false, nil
	}
	return false, fmt.Errorf("can't convert Float \"%v\" to bool", e.val)
}
