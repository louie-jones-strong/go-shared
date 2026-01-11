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
	for i, _ := range args {
		arg := args[i]

		// check if arg is comparable
		argType := reflect.TypeOf(arg)
		if !argType.Comparable() {
			return res, fmt.Errorf("arg %d type is not comparable: %s", i, argType.String())
		}

		res.args[i] = arg
	}

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
