package elements

import (
	"github.com/louie-jones-strong/go-shared/dataframe/apptype"
	"github.com/louie-jones-strong/go-shared/dataframe/series/elements/element"
)

type IElements interface {
	GetType() apptype.Type
	AllElems() []element.IElement
	Clone() IElements
	Elem(int) element.IElement
	Len() int
	Subset(indexes []int) IElements
	Append(values ...any)
}
