package routes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_parseURL(t *testing.T) {

	tests := []struct {
		name        string
		rawURL      string
		expectedRes string
		expectedErr error
	}{
		{
			name:        "empty URL",
			rawURL:      "",
			expectedRes: "",
			expectedErr: fmt.Errorf("parseURL called with empty URL"),
		},
		{
			name:        "valid URL",
			rawURL:      "https://example.com/path?query=123",
			expectedRes: "https://example.com/path?query=123",
			expectedErr: nil,
		},
		{
			name:        "URL with spaces",
			rawURL:      "https://example.com/Lean Diced Beef",
			expectedRes: "https://example.com/Lean%20Diced%20Beef",
			expectedErr: nil,
		},
		{
			name:        "URL with spaces and special characters",
			rawURL:      "https://example.com/3% Lean Diced Beef",
			expectedRes: "https://example.com/3%25%20Lean%20Diced%20Beef",
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			res, err := parseURL(tc.rawURL)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.Zero(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
			}
		})
	}
}
