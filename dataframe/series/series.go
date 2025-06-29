package series

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/stat"
)

type Type string

const (
	String Type = "string"
	Int    Type = "int"
	Float  Type = "float"
	Bool   Type = "bool"
)

type Element interface {
	Set(any)
	Val() any

	Type() Type

	ToString() string
	ToInt() (int, error)
	ToFloat() float64
	ToBool() (bool, error)
}
type Series struct {
	name     string
	t        Type
	elements IElements
}

func New(
	name string,
	t Type,
	values any,
) *Series {

	ret := &Series{
		name:     name,
		t:        t,
		elements: nil,
	}

	switch v := values.(type) {
	case []string:
		ret.elements = newElements(v, newStringElement)
	case []int:
		ret.elements = newElements(v, newIntElement)
	case []float64:
		ret.elements = newElements(v, newFloatElement)
	case []bool:
		ret.elements = newElements(v, newBoolElement)
	default:
		panic(fmt.Sprintf("unknown type %v", values))
	}

	return ret
}

func (s *Series) Len() int {
	return s.elements.Len()
}

func (s *Series) Append(values any) {

	news := New(s.name, s.t, values)
	switch s.t {
	case String:
		s.elements = append(s.elements.(Elements[*stringElement]), news.elements.(Elements[*stringElement])...)
	case Int:
		s.elements = append(s.elements.(Elements[*intElement]), news.elements.(Elements[*intElement])...)
	case Float:
		s.elements = append(s.elements.(Elements[*floatElement]), news.elements.(Elements[*floatElement])...)
	case Bool:
		s.elements = append(s.elements.(Elements[*boolElement]), news.elements.(Elements[*boolElement])...)
	}
}

func (s *Series) Val(i int) any {
	return s.Elem(i).Val()
}

func (s *Series) Elem(i int) Element {
	return s.elements.Elem(i)
}

func (s *Series) GetName() string {
	return s.name
}

func (s *Series) GetType() Type {
	return s.t
}

func (s Series) ToStrings() ([]string, error) {
	ret := make([]string, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.elements.Elem(i).ToString()
	}
	return ret, nil
}

func (s Series) ToInts() ([]int, error) {
	ret := make([]int, s.Len())
	for i := 0; i < s.Len(); i++ {
		val, err := s.elements.Elem(i).ToInt()
		if err != nil {
			return nil, err
		}
		ret[i] = val

	}
	return ret, nil
}

func (s Series) ToFloats() []float64 {
	ret := make([]float64, s.Len())
	for i := 0; i < s.Len(); i++ {
		e := s.elements.Elem(i)
		ret[i] = e.ToFloat()
	}
	return ret
}

func (s Series) ToBools() ([]bool, error) {
	ret := make([]bool, s.Len())
	for i := 0; i < s.Len(); i++ {
		val, err := s.elements.Elem(i).ToBool()
		if err != nil {
			return nil, err
		}
		ret[i] = val

	}
	return ret, nil
}

func (s Series) Sum() float64 {
	if s.Len() == 0 || s.GetType() == String || s.GetType() == Bool {
		return math.NaN()
	}
	sFloat := s.ToFloats()
	acc := float64(0)
	for i := 0; i < len(sFloat); i++ {
		acc += sFloat[i]
	}
	return acc
}

// StdDev calculates the standard deviation of a series
func (s Series) StdDev() float64 {
	stdDev := stat.StdDev(s.ToFloats(), nil)
	return stdDev
}

// Mean calculates the average value of a series
func (s Series) Mean() float64 {
	stdDev := stat.Mean(s.ToFloats(), nil)
	return stdDev
}

func (s Series) Min() float64 {
	if s.Len() == 0 || s.GetType() == String || s.GetType() == Bool {
		return math.NaN()
	}
	sFloat := s.ToFloats()

	minimum := sFloat[0]
	for i := 1; i < len(sFloat); i++ {
		minimum = min(minimum, sFloat[i])
	}
	return minimum
}

func (s Series) Max() float64 {
	if s.Len() == 0 || s.GetType() == String || s.GetType() == Bool {
		return math.NaN()
	}
	sFloat := s.ToFloats()

	maximum := sFloat[0]
	for i := 1; i < len(sFloat); i++ {
		maximum = max(maximum, sFloat[i])
	}
	return maximum
}

type Indexes interface {
}

func parseIndexes(l int, indexes Indexes) ([]int, error) {
	var ret []int
	switch tIndexes := indexes.(type) {
	case []int:
		ret = tIndexes
	case int:
		ret = []int{tIndexes}
	case []bool:
		bools := tIndexes
		if len(bools) != l {
			return nil, fmt.Errorf("indexing error: index dimensions mismatch")
		}
		for i, b := range bools {
			if b {
				ret = append(ret, i)
			}
		}
	default:
		return nil, fmt.Errorf("indexing error: unknown indexing mode")
	}
	return ret, nil
}

func (s Series) Subset(indexes Indexes) Series {
	idx, err := parseIndexes(s.Len(), indexes)
	if err != nil {
		return s
	}

	return Series{
		name:     s.name,
		t:        s.t,
		elements: s.elements.Subset(idx),
	}
}
