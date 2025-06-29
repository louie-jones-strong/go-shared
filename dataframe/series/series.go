package series

import (
	"fmt"
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
	// ToInt() (int, error)
	// ToFloat() float64
	// ToBool() (bool, error)
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
