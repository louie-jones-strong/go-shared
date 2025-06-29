package series

type stringElement struct {
	val string
}

func newStringElement(val string) *stringElement {
	return &stringElement{val: val}
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
