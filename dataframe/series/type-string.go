package series

var _ Element = (*stringElement)(nil)

type stringElement struct {
	e string
}

func (e stringElement) Type() Type {
	return String
}

func (e *stringElement) Set(value interface{}) {

	switch val := value.(type) {
	case string:
		e.e = string(val)
	default:
		return
	}
}
