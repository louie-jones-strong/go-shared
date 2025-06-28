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

type Elements interface {
	Elem(int) Element
	Len() int
}

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
	elements Elements
	t        Type
}

func New(
	name string,
	t Type,
	values any,
) *Series {
	ret := &Series{
		name: name,
		t:    t,
	}

	// Pre-allocate elements
	preAlloc := func(n int) {
		switch t {
		case String:
			ret.elements = make(stringElements, n)
		case Int:
			ret.elements = make(intElements, n)
		default:
			panic(fmt.Sprintf("unknown type %v", t))
		}
	}

	if values == nil {
		preAlloc(1)
		ret.elements.Elem(0).Set(nil)
		return ret
	}

	switch v := values.(type) {
	case []string:
		l := len(v)
		preAlloc(l)
		for i := 0; i < l; i++ {
			ret.elements.Elem(i).Set(v[i])
		}
	case []float64:
		l := len(v)
		preAlloc(l)
		for i := 0; i < l; i++ {
			ret.elements.Elem(i).Set(v[i])
		}
	case []int:
		l := len(v)
		preAlloc(l)
		for i := 0; i < l; i++ {
			ret.elements.Elem(i).Set(v[i])
		}
	case []bool:
		l := len(v)
		preAlloc(l)
		for i := 0; i < l; i++ {
			ret.elements.Elem(i).Set(v[i])
		}
	case Series:
		l := v.Len()
		preAlloc(l)
		for i := 0; i < l; i++ {
			ret.elements.Elem(i).Set(v.elements.Elem(i))
		}
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
		s.elements = append(s.elements.(stringElements), news.elements.(stringElements)...)
	case Int:
		s.elements = append(s.elements.(intElements), news.elements.(intElements)...)
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
