package element

import "github.com/louie-jones-strong/go-shared/dataframe/apptype"

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
	ToFloat() float64
	ToBool() (bool, error)
}
