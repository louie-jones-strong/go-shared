package elements

import (
	"fmt"

	"github.com/louie-jones-strong/go-shared/dataframe/apptype"
)

type IntElement struct {
	BaseElement[int]
}

func NewIntElement(val int) *IntElement {
	return &IntElement{
		NewBaseElement(val),
	}
}

func (e IntElement) Type() apptype.Type {
	return apptype.String
}

func (e IntElement) ToInt() (int, error) {
	return e.val, nil
}

func (e IntElement) ToFloat() float64 {
	return float64(e.val)
}
func (e IntElement) ToBool() (bool, error) {
	switch e.val {
	case 1:
		return true, nil
	case 0:
		return false, nil
	}
	return false, fmt.Errorf("can't convert Float \"%v\" to bool", e.val)
}
