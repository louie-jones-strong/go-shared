package series

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type stringElement struct {
	val string
}

func newStringElement(val string) *stringElement {
	return &stringElement{val: val}
}

func (e stringElement) Type() Type {
	return String
}

func (e stringElement) Val() any {
	return e.val
}

func (e *stringElement) Set(value any) {

	switch val := value.(type) {
	case string:
		e.val = string(val)
	default:
		return
	}
}

func (e stringElement) ToString() string {
	return e.val
}

func (e stringElement) ToInt() (int, error) {
	return strconv.Atoi(e.val)
}

func (e stringElement) ToFloat() float64 {
	f, err := strconv.ParseFloat(e.val, 64)
	if err != nil {
		return math.NaN()
	}
	return f
}

func (e stringElement) ToBool() (bool, error) {
	switch strings.ToLower(e.val) {
	case "true", "t", "1":
		return true, nil
	case "false", "f", "0":
		return false, nil
	}
	return false, fmt.Errorf("can't convert String \"%v\" to bool", e.val)
}
