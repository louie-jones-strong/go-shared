package series

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/louie-jones-strong/go-shared/dataframe/apptype"
	"github.com/louie-jones-strong/go-shared/dataframe/series/elements"
	"github.com/louie-jones-strong/go-shared/dataframe/series/elements/element"
	"gonum.org/v1/gonum/stat"
)

type Series struct {
	name string
	elms elements.IElements
}

func BuildSeries(name string, values any) *Series {

	switch vals := values.(type) {
	case []string:
		return New(name, elements.BuildElements(vals, element.NewStringElement))
	case []int:
		return New(name, elements.BuildElements(vals, element.NewFloatElementFromInt))
	case []float64:
		return New(name, elements.BuildElements(vals, element.NewFloatElement))
	case []bool:
		return New(name, elements.BuildElements(vals, element.NewBoolElement))
	case []time.Time:
		return New(name, elements.BuildElements(vals, element.NewDateTimeElement))
	default:
		panic(fmt.Sprintf("unknown type %v", values))
	}
}

func New(name string, elms elements.IElements) *Series {
	return &Series{
		name: name,
		elms: elms,
	}
}

func (s *Series) Clone() *Series {
	return New(
		s.name,
		s.elms.Clone(),
	)
}

func (s *Series) GetType() apptype.Type {
	return s.elms.GetType()
}

func (s *Series) Rename(newName string) {
	s.name = newName
}

func (s *Series) Len() int {
	return s.elms.Len()
}

func (s *Series) Append(values any) {
	s.elms.Append(values)
}

func (s *Series) Val(i int) any {
	return s.Elem(i).Val()
}

func (s *Series) Elem(i int) element.IElement {
	return s.elms.Elem(i)
}

func (s *Series) GetName() string {
	return s.name
}

func (s Series) Values() []any {
	ret := make([]any, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.elms.Elem(i).Val()
	}
	return ret
}

func (s Series) ToFloats() []float64 {
	ret := make([]float64, s.Len())
	for i := 0; i < s.Len(); i++ {
		e := s.elms.Elem(i)
		ret[i] = e.ToFloat()
	}
	return ret
}

func (s Series) Sum() float64 {
	if s.Len() == 0 || s.GetType() == apptype.String || s.GetType() == apptype.Bool {
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
	if s.Len() == 0 || s.GetType() == apptype.String || s.GetType() == apptype.Bool {
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
	if s.Len() == 0 || s.GetType() == apptype.String || s.GetType() == apptype.Bool {
		return math.NaN()
	}
	sFloat := s.ToFloats()

	maximum := sFloat[0]
	for i := 1; i < len(sFloat); i++ {
		maximum = max(maximum, sFloat[i])
	}
	return maximum
}

type applyFunc[In element.IElement, Out element.IElement] func(item In) (Out, error)

func Apply[In element.IElement, Out element.IElement](
	elems elements.Elements[In],
	delegate applyFunc[In, Out],
) (elements.Elements[Out], error) {

	newElements := make([]Out, len(elems))
	for i := 0; i < elems.Len(); i++ {

		newElem, err := delegate(elems[i])
		if err != nil {
			return nil, err
		}

		newElements[i] = newElem

	}
	return newElements, nil
}

func (s *Series) ApplyInPlace(delegate applyFunc[element.IElement, element.IElement]) {
	newElms, err := Apply(s.elms.AllElems(), delegate)
	if err != nil {
		panic(err)
	}

	temp := elements.NewElements(newElms)
	s.elms = &temp
}

type Indexes interface{}

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
		name: s.name,
		elms: s.elms.Subset(idx),
	}
}

func (s Series) Order(reverse bool) []int {
	var ie indexedElements
	var nasIdx []int
	for i := 0; i < s.Len(); i++ {
		e := s.elms.Elem(i)
		ie = append(ie, indexedElement{i, e})
	}
	var srt sort.Interface
	srt = ie
	if reverse {
		srt = sort.Reverse(srt)
	}
	sort.Stable(srt)
	var ret []int
	for _, e := range ie {
		ret = append(ret, e.index)
	}
	return append(ret, nasIdx...)
}

type indexedElement struct {
	index   int
	element element.IElement
}

type indexedElements []indexedElement

func (e indexedElements) Len() int           { return len(e) }
func (e indexedElements) Less(i, j int) bool { return e[i].element.Less(e[j].element) }
func (e indexedElements) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
