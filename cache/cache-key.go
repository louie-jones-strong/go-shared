package cache

import (
	"fmt"
	"reflect"
	"runtime"
)

type cacheKey struct {
	funcName string
	args     [10]any
}

func newCacheKey(f any, args ...any) (cacheKey, error) {
	res := cacheKey{}

	if len(args) > 10 {
		return res, fmt.Errorf("to many args")
	}

	res.args = [10]any{}
	copy(res.args[:], args)

	res.funcName = getFunctionName(f)

	return res, nil
}

func (ck *cacheKey) ToString() string {
	return ck.funcName + fmt.Sprintf(" Args: %v", ck.args)
}

func getFunctionName(f any) string {
	fn := runtime.FuncForPC(reflect.ValueOf(f).Pointer())

	fullName := fn.Name()

	return fullName
}
