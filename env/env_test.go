package env

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_GetKey(t *testing.T) {

	tests := []struct {
		name        string
		key         string
		expectedRes string
		expectedErr error
	}{
		{
			name:        "empty key",
			key:         "",
			expectedRes: "",
			expectedErr: fmt.Errorf("environment variable \"%v\" not set. Please add it to your .env file", ""),
		},
		{
			name:        "missing key",
			key:         "missing_key",
			expectedRes: "",
			expectedErr: fmt.Errorf("environment variable \"%v\" not set. Please add it to your .env file", "missing_key"),
		},
		{
			name:        "found key",
			key:         "found_key",
			expectedRes: "found_value",
			expectedErr: nil,
		},
	}

	err := LoadEnv()
	assert.NoError(t, err)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			res, err := GetKey(tc.key)

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
