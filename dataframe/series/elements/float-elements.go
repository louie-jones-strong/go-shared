package elements

import (
	"fmt"

	"github.com/louie-jones-strong/go-shared/dataframe/apptype"
)

type FloatElement struct {
	BaseElement[float64]
}

func NewFloatElement(val float64) *FloatElement {
	return &FloatElement{
		NewBaseElement(val),
	}
}

func (e FloatElement) Type() apptype.Type {
	return apptype.String
}

func (e FloatElement) ToInt() (int, error) {
	return int(e.val), nil
}

func (e FloatElement) ToFloat() float64 {
	return e.val
}

func (e FloatElement) ToBool() (bool, error) {
	switch e.val {
	case 1:
		return true, nil
	case 0:
		return false, nil
	}
	return false, fmt.Errorf("can't convert Float \"%v\" to bool", e.val)
}
