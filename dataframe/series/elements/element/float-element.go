package element

import (
	"math"

	"github.com/louie-jones-strong/go-shared/dataframe/apptype"
)

type FloatElement struct {
	BaseElement[float64]
}

func NewFloatElementFromInt(val int) *FloatElement {
	return NewFloatElement(float64(val))
}
func NewFloatElement(val float64) *FloatElement {
	return &FloatElement{
		NewBaseElement(val),
	}
}

func (e *FloatElement) Clone() IElement {
	return &FloatElement{
		e.BaseElement.Clone(),
	}
}

func (e FloatElement) Type() apptype.Type {
	return apptype.String
}

func (e FloatElement) ToFloat() float64 {
	return e.val
}

func (e FloatElement) Less(elem IElement) bool {
	f := elem.ToFloat()
	if math.IsNaN(f) {
		return false
	}
	return e.val < f
}

func (e FloatElement) LessEq(elem IElement) bool {
	f := elem.ToFloat()
	if math.IsNaN(f) {
		return false
	}
	return e.val <= f
}

func (e FloatElement) Greater(elem IElement) bool {
	f := elem.ToFloat()
	if math.IsNaN(f) {
		return false
	}
	return e.val > f
}

func (e FloatElement) GreaterEq(elem IElement) bool {
	f := elem.ToFloat()
	if math.IsNaN(f) {
		return false
	}
	return e.val >= f
}
