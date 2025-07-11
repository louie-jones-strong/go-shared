package elements

import (
	"fmt"
	"time"

	"github.com/louie-jones-strong/go-shared/dataframe/apptype"
)

type DateTimeElement struct {
	BaseElement[time.Time]
}

func NewDateTimeElement(val time.Time) *DateTimeElement {
	return &DateTimeElement{
		NewBaseElement(val),
	}
}

func (e DateTimeElement) Clone() IElement {
	return &DateTimeElement{
		e.BaseElement.Clone(),
	}
}

func (e DateTimeElement) Type() apptype.Type {
	return apptype.String
}

func (e DateTimeElement) ToString() string {
	return e.val.Format(time.RFC3339)
}

func (e DateTimeElement) ToInt() (int, error) {
	return int(e.val.Unix()), nil
}

func (e DateTimeElement) ToFloat() float64 {
	return float64(e.val.Unix())
}

func (e DateTimeElement) ToBool() (bool, error) {
	return false, fmt.Errorf("can't convert datetime \"%v\" to bool", e.val)
}

func (e DateTimeElement) Less(elem IElement) bool {
	dt, ok := elem.(*DateTimeElement)
	if !ok {
		return false
	}
	return e.val.Before(dt.val)
}

func (e DateTimeElement) LessEq(elem IElement) bool {
	dt, ok := elem.(*DateTimeElement)
	if !ok {
		return false
	}
	return e.val.Before(dt.val) || e.val.Equal(dt.val)
}

func (e DateTimeElement) Greater(elem IElement) bool {
	dt, ok := elem.(*DateTimeElement)
	if !ok {
		return false
	}
	return e.val.After(dt.val)
}

func (e DateTimeElement) GreaterEq(elem IElement) bool {
	dt, ok := elem.(*DateTimeElement)
	if !ok {
		return false
	}
	return e.val.After(dt.val) || e.val.Equal(dt.val)
}
