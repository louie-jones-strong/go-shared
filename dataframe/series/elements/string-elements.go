package elements

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/louie-jones-strong/go-shared/dataframe/apptype"
)

type StringElement struct {
	BaseElement[string]
}

func NewStringElement(val string) *StringElement {
	return &StringElement{
		BaseElement: NewBaseElement(val),
	}
}

func (e StringElement) Type() apptype.Type {
	return apptype.String
}

func (e StringElement) ToInt() (int, error) {
	return strconv.Atoi(e.val)
}

func (e StringElement) ToFloat() float64 {
	f, err := strconv.ParseFloat(e.val, 64)
	if err != nil {
		return math.NaN()
	}
	return f
}

func (e StringElement) ToBool() (bool, error) {
	switch strings.ToLower(e.val) {
	case "true", "t", "1":
		return true, nil
	case "false", "f", "0":
		return false, nil
	}
	return false, fmt.Errorf("can't convert String \"%v\" to bool", e.val)
}
