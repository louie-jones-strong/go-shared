package series

type IElements interface {
	Elem(int) Element
	Len() int
	Subset(indexes []int) IElements
}

type Elements[T Element] []T

func newElements[V any, E Element](
	values []V,
	elemBuilder func(V) E,
) Elements[E] {
	elements := make(Elements[E], len(values))
	for i := 0; i < len(values); i++ {
		elements[i] = elemBuilder(values[i])
	}

	return elements
}

func (e Elements[T]) Len() int           { return len(e) }
func (e Elements[T]) Elem(i int) Element { return e[i] }

func (e Elements[T]) Subset(indexes []int) IElements {

	ret := make(Elements[T], len(indexes))
	for k, i := range indexes {
		ret[k] = e[i]
	}
	return ret
}
