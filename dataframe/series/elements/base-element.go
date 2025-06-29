package elements

import "fmt"

type BaseElement[T any] struct {
	val T
}

func NewBaseElement[T any](val T) BaseElement[T] {
	return BaseElement[T]{val: val}
}

func (e BaseElement[T]) Val() any {
	return e.val
}

func (e BaseElement[T]) Set(value any) {

	switch val := value.(type) {
	case T:
		e.val = val
	default:
		return
	}
}

func (e BaseElement[T]) ToString() string {
	return fmt.Sprintf("%v", e.val)
}
