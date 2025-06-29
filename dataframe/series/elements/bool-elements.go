package elements

import "github.com/louie-jones-strong/go-shared/dataframe/apptype"

type BoolElement struct {
	BaseElement[bool]
}

func NewBoolElement(val bool) *BoolElement {
	return &BoolElement{
		NewBaseElement(val),
	}
}

func (e BoolElement) Type() apptype.Type {
	return apptype.String
}

func (e BoolElement) ToInt() (int, error) {
	if e.val {
		return 1, nil
	}
	return 0, nil
}

func (e BoolElement) ToFloat() float64 {
	if e.val {
		return 1.0
	}
	return 0.0
}

func (e BoolElement) ToBool() (bool, error) {
	return e.val, nil
}
