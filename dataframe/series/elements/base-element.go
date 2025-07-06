package elements

import (
	"fmt"
	"time"
)

type supportedTypes interface {
	string | int | float64 | bool | time.Time
}

type BaseElement[T supportedTypes] struct {
	val T
}

func NewBaseElement[T supportedTypes](val T) BaseElement[T] {
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

func (e BaseElement[T]) Eq(other Element) bool {
	otherVal, ok := other.Val().(T)
	if !ok {
		return false
	}

	return e.val == otherVal
}

func (e BaseElement[T]) Neq(other Element) bool {
	otherVal, ok := other.Val().(T)
	if !ok {
		return true
	}

	return e.val != otherVal
}
