package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_ConvertMapToKVPList(t *testing.T) {

	tests := []struct {
		name        string
		dict        map[string]int
		expectedRes []KVP[string, int]
	}{
		{
			name:        "nil map",
			dict:        nil,
			expectedRes: []KVP[string, int]{},
		},
		{
			name:        "empty map",
			dict:        map[string]int{},
			expectedRes: []KVP[string, int]{},
		},
		{
			name:        "1 item map",
			dict:        map[string]int{"a": 1},
			expectedRes: []KVP[string, int]{{key: "a", value: 1}},
		},
		{
			name:        "multiple items map",
			dict:        map[string]int{"a": 1, "b": 2, "c": 3},
			expectedRes: []KVP[string, int]{{key: "a", value: 1}, {key: "b", value: 2}, {key: "c", value: 3}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			res := ConvertMapToKVPList(tc.dict)

			assert.ElementsMatch(t, tc.expectedRes, res)
		})
	}
}

func TestUnit_GetMapKeys(t *testing.T) {

	tests := []struct {
		name        string
		dict        map[string]int
		expectedRes []string
	}{
		{
			name:        "nil map",
			dict:        nil,
			expectedRes: []string{},
		},
		{
			name:        "empty map",
			dict:        map[string]int{},
			expectedRes: []string{},
		},
		{
			name:        "1 item map",
			dict:        map[string]int{"a": 1},
			expectedRes: []string{"a"},
		},
		{
			name:        "multiple items map",
			dict:        map[string]int{"a": 1, "b": 2, "c": 3},
			expectedRes: []string{"a", "b", "c"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			res := GetMapKeys(tc.dict)

			assert.ElementsMatch(t, tc.expectedRes, res)
		})
	}
}

func TestUnit_GetMapValues(t *testing.T) {

	tests := []struct {
		name        string
		dict        map[string]int
		expectedRes []int
	}{
		{
			name:        "nil map",
			dict:        nil,
			expectedRes: []int{},
		},
		{
			name:        "empty map",
			dict:        map[string]int{},
			expectedRes: []int{},
		},
		{
			name:        "1 item map",
			dict:        map[string]int{"a": 1},
			expectedRes: []int{1},
		},
		{
			name:        "multiple items map",
			dict:        map[string]int{"a": 1, "b": 2, "c": 3},
			expectedRes: []int{1, 2, 3},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			res := GetMapValues(tc.dict)

			assert.ElementsMatch(t, tc.expectedRes, res)
		})
	}
}
