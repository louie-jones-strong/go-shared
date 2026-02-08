package routes

import (
	"fmt"
	"math"
)

type Page[T any] []T

func SplitCollectionIntoPages[T any](items []T, itemsPerPage int) []Page[T] {

	numItems := len(items)
	numPages := int(math.Ceil(float64(numItems) / float64(itemsPerPage)))

	pages := make([]Page[T], 0, numPages)
	for i := 0; i < numPages; i++ {

		startIdx := i * itemsPerPage
		endIdx := startIdx + itemsPerPage
		if endIdx > numItems {
			endIdx = numItems
		}
		pages = append(pages, Page[T](items[startIdx:endIdx]))
	}

	return pages
}

type PaginatedResponse[T any] struct {
	Items        Page[T] `json:"items"`
	CurrentPage  int     `json:"current-page-index"`
	PreviousPage int     `json:"previous-page-index,omitempty"`
	NextPage     int     `json:"next-page-index,omitempty"`
}

func GetPaginatedResponse[T any](allItems []T, currentPageIdx int, itemsPerPage int) (PaginatedResponse[T], error) {
	pages := SplitCollectionIntoPages(allItems, itemsPerPage)

	if currentPageIdx < 0 || currentPageIdx >= len(pages) {
		return PaginatedResponse[T]{}, fmt.Errorf("invalid page index %d should be in range 0 to %d", currentPageIdx, len(pages)-1)
	}

	prevPageIdx := currentPageIdx - 1
	if prevPageIdx < 0 {
		prevPageIdx = -1
	}
	nextPageIdx := currentPageIdx + 1
	if nextPageIdx >= len(pages) {
		nextPageIdx = -1
	}

	return PaginatedResponse[T]{
		Items:        pages[currentPageIdx],
		CurrentPage:  currentPageIdx,
		PreviousPage: prevPageIdx,
		NextPage:     nextPageIdx,
	}, nil
}
