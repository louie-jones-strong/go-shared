package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_CacheCall(t *testing.T) {

	t.Run("no sub caches", func(t *testing.T) {

		idx := 0

		tempFunc := func() (int, error) {
			idx++
			return idx, nil
		}

		res, err := tempFunc()
		assert.NoError(t, err)
		assert.Equal(t, 1, res)

		res, err = tempFunc()
		assert.NoError(t, err)
		assert.Equal(t, 2, res)

		res, err = CacheCall(tempFunc)
		assert.NoError(t, err)
		assert.Equal(t, 3, res)

		res, err = CacheCall(tempFunc)
		assert.NoError(t, err)
		assert.Equal(t, 3, res)
	})

	t.Run("with sub caches", func(t *testing.T) {

		idx := 0

		func1 := func() (int, error) {
			idx++
			return idx, nil
		}
		func2 := func() (int, error) {
			return CacheCall(func1)
		}

		res, err := CacheCall(func2)
		assert.NoError(t, err)
		assert.Equal(t, 1, res)

		res, err = CacheCall(func2)
		assert.NoError(t, err)
		assert.Equal(t, 1, res)

		res, err = CacheCall(func1)
		assert.NoError(t, err)
		assert.Equal(t, 1, res)
	})
}
