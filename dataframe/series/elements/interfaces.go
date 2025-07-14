package elements

import "github.com/louie-jones-strong/go-shared/dataframe/apptype"

type IElements interface {
	GetType() apptype.Type
	AllElems() []IElement
	Clone() IElements
	Elem(int) IElement
	Len() int
	Subset(indexes []int) IElements
	Append(values ...any)
}

type IElement interface {
	Clone() IElement

	Set(any)

	Eq(IElement) bool
	Neq(IElement) bool
	Less(IElement) bool
	LessEq(IElement) bool
	Greater(IElement) bool
	GreaterEq(IElement) bool

	Val() any

	Type() apptype.Type

	ToString() string
	ToInt() (int, error)
	ToFloat() float64
	ToBool() (bool, error)
}
