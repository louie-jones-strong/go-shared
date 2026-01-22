package element

import "github.com/louie-jones-strong/go-shared/collections/dataframe/apptype"

type BoolElement struct {
	BaseElement[bool]
}

func NewBoolElement(val bool) *BoolElement {
	return &BoolElement{
		NewBaseElement(val),
	}
}

func (e *BoolElement) Clone() IElement {
	return &BoolElement{
		e.BaseElement.Clone(),
	}
}

func (e BoolElement) Type() apptype.Type {
	return apptype.String
}

func (e BoolElement) ToFloat() float64 {
	if e.val {
		return 1.0
	}
	return 0.0
}

func (e BoolElement) Less(elem IElement) bool {
	b, ok := elem.(*BoolElement)
	if !ok {
		return false
	}
	return !e.val && b.val
}

func (e BoolElement) LessEq(elem IElement) bool {
	b, ok := elem.(*BoolElement)
	if !ok {
		return false
	}
	return !e.val || b.val
}

func (e BoolElement) Greater(elem IElement) bool {
	b, ok := elem.(*BoolElement)
	if !ok {
		return false
	}
	return e.val && !b.val
}

func (e BoolElement) GreaterEq(elem IElement) bool {
	b, ok := elem.(*BoolElement)
	if !ok {
		return false
	}
	return e.val || !b.val
}
