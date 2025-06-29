package series

type boolElement struct {
	val bool
}

func newBoolElement(val bool) *boolElement {
	return &boolElement{val: val}
}

func (e boolElement) Type() Type {
	return String
}

func (e boolElement) Val() any {
	return e.val
}

func (e *boolElement) Set(value any) {

	switch val := value.(type) {
	case bool:
		e.val = bool(val)
	default:
		return
	}
}

func (e boolElement) ToString() string {
	if e.val {
		return "true"
	}
	return "false"
}
