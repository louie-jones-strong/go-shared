package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/louie-jones-strong/go-shared/cache"
	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	Value string
	Count int
}

func TestUnit_CachedFileStorage_IsValid(t *testing.T) {

	t.Run("initially false when not loaded", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		err := os.WriteFile(filePath, []byte{}, 0644)
		assert.NoError(t, err)

		mockStore := &MockStorage[TestModel]{}
		cached := NewCachedFileStorage(filePath, mockStore)

		assert.False(t, cached.IsValid())
	})

	t.Run("true after successful load", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		err := os.WriteFile(filePath, []byte{}, 0644)
		assert.NoError(t, err)

		mockStore := &MockStorage[TestModel]{
			Data: TestModel{Value: "test", Count: 42},
		}
		cached := NewCachedFileStorage(filePath, mockStore)

		_, err = cached.Load()
		assert.NoError(t, err)

		assert.True(t, cached.IsValid())
	})

	t.Run("false after file modification", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		err := os.WriteFile(filePath, []byte{}, 0644)
		assert.NoError(t, err)

		mockStore := &MockStorage[TestModel]{
			Data: TestModel{Value: "test", Count: 42},
		}
		cached := NewCachedFileStorage(filePath, mockStore)

		_, err = cached.Load()
		assert.NoError(t, err)
		assert.True(t, cached.IsValid())

		time.Sleep(10 * time.Millisecond)
		err = os.WriteFile(filePath, []byte("modified"), 0644)
		assert.NoError(t, err)

		assert.False(t, cached.IsValid())
		assert.Nil(t, cached.lastModTime)
	})

	t.Run("false when file does not exist", func(t *testing.T) {
		filePath := "/nonexistent/path/file.txt"

		mockStore := &MockStorage[TestModel]{}
		cached := NewCachedFileStorage(filePath, mockStore)

		assert.False(t, cached.IsValid())
	})
}

func TestUnit_CachedFileStorage_Load(t *testing.T) {

	t.Run("first load from file", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		err := os.WriteFile(filePath, []byte{}, 0644)
		assert.NoError(t, err)

		expectedData := TestModel{Value: "test", Count: 42}
		mockStore := &MockStorage[TestModel]{
			Data: expectedData,
		}
		cached := NewCachedFileStorage(filePath, mockStore)

		result, err := cached.Load()

		assert.NoError(t, err)
		assert.Equal(t, expectedData, result)
		assert.NotNil(t, cached.lastModTime)
		assert.Equal(t, expectedData, cached.cachedContent)
	})

	t.Run("cache hit returns cached data", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		err := os.WriteFile(filePath, []byte{}, 0644)
		assert.NoError(t, err)

		expectedData := TestModel{Value: "test", Count: 42}
		mockStore := &MockStorage[TestModel]{
			Data: expectedData,
		}
		cached := NewCachedFileStorage(filePath, mockStore)

		result1, err := cached.Load()
		assert.NoError(t, err)
		assert.Equal(t, expectedData, result1)

		mockStore.Data = TestModel{Value: "different", Count: 99}

		result2, err := cached.Load()
		assert.NoError(t, err)
		assert.Equal(t, expectedData, result2)
	})

	t.Run("cache miss after file modification", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		err := os.WriteFile(filePath, []byte{}, 0644)
		assert.NoError(t, err)

		mockStore := &MockStorage[TestModel]{
			Data: TestModel{Value: "initial", Count: 1},
		}
		cached := NewCachedFileStorage(filePath, mockStore)

		result1, err := cached.Load()
		assert.NoError(t, err)
		assert.Equal(t, "initial", result1.Value)

		time.Sleep(10 * time.Millisecond)
		err = os.WriteFile(filePath, []byte("modified"), 0644)
		assert.NoError(t, err)

		newData := TestModel{Value: "updated", Count: 99}
		mockStore.Data = newData

		result2, err := cached.Load()
		assert.NoError(t, err)
		assert.Equal(t, newData, result2)
		assert.Equal(t, newData, cached.cachedContent)
	})

	t.Run("error when file not found", func(t *testing.T) {
		filePath := "/nonexistent/path/file.txt"

		mockStore := &MockStorage[TestModel]{
			Data: TestModel{Value: "test", Count: 42},
		}
		cached := NewCachedFileStorage(filePath, mockStore)

		_, err := cached.Load()

		assert.Error(t, err)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("error from sub store", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		err := os.WriteFile(filePath, []byte{}, 0644)
		assert.NoError(t, err)

		mockStore := &MockStorage[TestModel]{
			LoadError: assert.AnError,
		}
		cached := NewCachedFileStorage(filePath, mockStore)

		_, err = cached.Load()

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})
}

func TestUnit_CachedFileStorage_Save(t *testing.T) {

	t.Run("invalidates cache", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		err := os.WriteFile(filePath, []byte{}, 0644)
		assert.NoError(t, err)

		mockStore := &MockStorage[TestModel]{
			Data: TestModel{Value: "initial", Count: 1},
		}
		cached := NewCachedFileStorage(filePath, mockStore)

		_, err = cached.Load()
		assert.NoError(t, err)
		assert.True(t, cached.IsValid())

		newData := TestModel{Value: "saved", Count: 2}
		err = cached.Save(newData)

		assert.NoError(t, err)
		assert.False(t, cached.IsValid())
		assert.Nil(t, cached.lastModTime)
		assert.Equal(t, TestModel{}, cached.cachedContent)
	})

	t.Run("returns error from sub store", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		err := os.WriteFile(filePath, []byte{}, 0644)
		assert.NoError(t, err)

		mockStore := &MockStorage[TestModel]{
			SaveError: assert.AnError,
		}
		cached := NewCachedFileStorage(filePath, mockStore)

		err = cached.Save(TestModel{Value: "test", Count: 42})

		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})
}

func TestUnit_CachedFileStorage_WithFuncCache(t *testing.T) {

	t.Run("basic integration", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		err := os.WriteFile(filePath, []byte{}, 0644)
		assert.NoError(t, err)

		mockStore := &MockStorage[TestModel]{
			Data: TestModel{Value: "test", Count: 42},
		}
		cached := NewCachedFileStorage(filePath, mockStore)

		loadFunc := func() (TestModel, error) {
			return cached.Load()
		}

		result1, err := cache.CacheCall(loadFunc)
		assert.NoError(t, err)
		assert.Equal(t, "test", result1.Value)
		assert.Equal(t, 42, result1.Count)

		mockStore.Data = TestModel{Value: "changed", Count: 99}

		result2, err := cache.CacheCall(loadFunc)
		assert.NoError(t, err)
		assert.Equal(t, "test", result2.Value)
		assert.Equal(t, 42, result2.Count)
	})

	t.Run("nested caching with sub scopes", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		err := os.WriteFile(filePath, []byte{}, 0644)
		assert.NoError(t, err)

		mockStore := &MockStorage[TestModel]{
			Data: TestModel{Value: "nested", Count: 10},
		}
		cached := NewCachedFileStorage(filePath, mockStore)

		innerFunc := func() (TestModel, error) {
			return cached.Load()
		}

		outerFunc := func() (string, error) {
			data, err := cache.CacheCall(innerFunc)
			if err != nil {
				return "", err
			}
			return data.Value + "_processed", nil
		}

		result1, err := cache.CacheCall(outerFunc)
		assert.NoError(t, err)
		assert.Equal(t, "nested_processed", result1)

		mockStore.Data = TestModel{Value: "different", Count: 20}

		result2, err := cache.CacheCall(outerFunc)
		assert.NoError(t, err)
		assert.Equal(t, "nested_processed", result2)
	})

	t.Run("multiple files with independent caches", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath1 := filepath.Join(tmpDir, "test1.txt")
		filePath2 := filepath.Join(tmpDir, "test2.txt")
		err := os.WriteFile(filePath1, []byte{}, 0644)
		assert.NoError(t, err)
		err = os.WriteFile(filePath2, []byte{}, 0644)
		assert.NoError(t, err)

		mockStore1 := &MockStorage[TestModel]{
			Data: TestModel{Value: "file1", Count: 1},
		}
		cached1 := NewCachedFileStorage(filePath1, mockStore1)

		mockStore2 := &MockStorage[TestModel]{
			Data: TestModel{Value: "file2", Count: 2},
		}
		cached2 := NewCachedFileStorage(filePath2, mockStore2)

		loadFunc1 := func() (TestModel, error) {
			return cached1.Load()
		}

		loadFunc2 := func() (TestModel, error) {
			return cached2.Load()
		}

		result1, err := cache.CacheCall(loadFunc1)
		assert.NoError(t, err)
		assert.Equal(t, "file1", result1.Value)

		result2, err := cache.CacheCall(loadFunc2)
		assert.NoError(t, err)
		assert.Equal(t, "file2", result2.Value)

		time.Sleep(10 * time.Millisecond)
		err = os.WriteFile(filePath1, []byte("modified"), 0644)
		assert.NoError(t, err)

		mockStore1.Data = TestModel{Value: "updated1", Count: 10}
		mockStore2.Data = TestModel{Value: "updated2", Count: 20}

		result1Again, err := cache.CacheCall(loadFunc1)
		assert.NoError(t, err)
		assert.Equal(t, "updated1", result1Again.Value)

		result2Again, err := cache.CacheCall(loadFunc2)
		assert.NoError(t, err)
		assert.Equal(t, "file2", result2Again.Value)
	})

	t.Run("file modification invalidates full tree", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		err := os.WriteFile(filePath, []byte{}, 0644)
		assert.NoError(t, err)

		mockStore := &MockStorage[TestModel]{
			Data: TestModel{Value: "initial", Count: 1},
		}
		cached := NewCachedFileStorage(filePath, mockStore)

		loadFunc := func() (TestModel, error) {
			return cached.Load()
		}

		result1, err := cache.CacheCall(loadFunc)
		assert.NoError(t, err)
		assert.Equal(t, "initial", result1.Value)

		time.Sleep(10 * time.Millisecond)
		err = os.WriteFile(filePath, []byte("modified"), 0644)
		assert.NoError(t, err)

		mockStore.Data = TestModel{Value: "updated", Count: 2}

		result2, err := cache.CacheCall(loadFunc)
		assert.NoError(t, err)
		assert.Equal(t, "updated", result2.Value)
	})

	t.Run("save does not invalidate FuncCache", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		err := os.WriteFile(filePath, []byte{}, 0644)
		assert.NoError(t, err)

		mockStore := &MockStorage[TestModel]{
			Data: TestModel{Value: "initial", Count: 1},
		}
		cached := NewCachedFileStorage(filePath, mockStore)

		loadFunc := func() (TestModel, error) {
			return cached.Load()
		}

		result1, err := cache.CacheCall(loadFunc)
		assert.NoError(t, err)
		assert.Equal(t, "initial", result1.Value)

		err = cached.Save(TestModel{Value: "saved", Count: 5})
		assert.NoError(t, err)
		assert.False(t, cached.IsValid())

		result2, err := cache.CacheCall(loadFunc)
		assert.NoError(t, err)
		assert.Equal(t, "initial", result2.Value)
	})
}
