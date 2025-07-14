package elements

import (
	"github.com/louie-jones-strong/go-shared/dataframe/apptype"
	"github.com/louie-jones-strong/go-shared/dataframe/series/elements/element"
)

type Elements[T element.IElement] []T

func NewElements[E element.IElement](items []E) Elements[E] {
	return items
}

func BuildElements[V any, E element.IElement](
	values []V,
	elemBuilder func(V) E,
) Elements[E] {
	elements := make(Elements[E], len(values))
	for i := 0; i < len(values); i++ {
		elements[i] = elemBuilder(values[i])
	}

	return elements
}

func (e Elements[T]) GetType() apptype.Type {
	return getType(e)
}

func getType(e IElements) apptype.Type {
	switch e.(type) {
	case Elements[*element.StringElement]:
		return apptype.String
	case Elements[*element.IntElement]:
		return apptype.Int
	case Elements[*element.FloatElement]:
		return apptype.Float
	case Elements[*element.BoolElement]:
		return apptype.Bool
	case Elements[*element.DateTimeElement]:
		return apptype.DateTime
	default:
		return apptype.None
	}
}

func (e Elements[T]) AllElems() []element.IElement {
	ret := make([]element.IElement, len(e))
	for i := 0; i < len(e); i++ {
		ret[i] = e[i]
	}
	return ret
}

func (e Elements[T]) Clone() IElements {
	newElems := make([]T, len(e))
	for i := 0; i < len(e); i++ {

		cloned := e[i].Clone()
		typedCloned, ok := cloned.(T)
		if !ok {
			panic("Clone() called on element but it returned wrong type")
		}
		newElems[i] = typedCloned
	}

	return NewElements(newElems)
}

func (e Elements[T]) Len() int {
	return len(e)
}

func (e Elements[T]) Elem(i int) element.IElement {
	return e[i]
}

func (e Elements[T]) Subset(indexes []int) IElements {
	ret := make(Elements[T], len(indexes))
	for k, i := range indexes {
		ret[k] = e[i]
	}
	return ret
}

func (e Elements[T]) Append(values ...any) {

	for _, value := range values {
		val, ok := value.(T)
		if !ok {
			continue
		}

		e = append(e, val)
	}

}
