package elements

type Elements[T IElement] []T

func NewElements[E IElement](items []E) Elements[E] {
	return items
}

func BuildElements[V any, E IElement](
	values []V,
	elemBuilder func(V) E,
) Elements[E] {
	elements := make(Elements[E], len(values))
	for i := 0; i < len(values); i++ {
		elements[i] = elemBuilder(values[i])
	}

	return elements
}

func (e Elements[T]) AllElems() []IElement {
	ret := make([]IElement, len(e))
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

func (e Elements[T]) Elem(i int) IElement {
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
