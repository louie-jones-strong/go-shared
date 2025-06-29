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

func (e boolElement) ToInt() (int, error) {
	if e.val {
		return 1, nil
	}
	return 0, nil
}

func (e boolElement) ToFloat() float64 {
	if e.val {
		return 1.0
	}
	return 0.0
}

func (e boolElement) ToBool() (bool, error) {
	return e.val, nil
}
