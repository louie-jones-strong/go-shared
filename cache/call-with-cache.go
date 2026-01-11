package cache

func CacheCall[V any](f func() (V, error)) (V, error) {
	var defaultOut V

	ck, err := newCacheKey(f)
	if err != nil {
		return defaultOut, err
	}

	cs := GetCacheService()
	ci := FindOrNewFuncCache(cs, ck)
	cs.AddOpenScope(ci)
	defer cs.CloseScope(ci)

	if ci.IsValid() {
		return GetVal[V](ci)
	}

	val, err := f()
	if err != nil {
		return defaultOut, err
	}

	ci.SetValue(val)
	return val, nil
}

func CacheCall1Arg[
	a1 any,
	V any,
](f func(a1) (V, error), arg1 a1) (V, error) {
	var defaultOut V

	ck, err := newCacheKey(f, arg1)
	if err != nil {
		return defaultOut, err
	}

	cs := GetCacheService()
	ci := FindOrNewFuncCache(cs, ck)
	cs.AddOpenScope(ci)
	defer cs.CloseScope(ci)

	if ci.IsValid() {
		return GetVal[V](ci)
	}

	val, err := f(arg1)
	if err != nil {
		return defaultOut, err
	}

	ci.SetValue(val)
	return val, nil
}

func CacheCall2Args[
	a1 any,
	a2 any,
	V any,
](f func(a1, a2) (V, error), arg1 a1, arg2 a2) (V, error) {
	var defaultOut V

	ck, err := newCacheKey(f, arg1, arg2)
	if err != nil {
		return defaultOut, err
	}

	cs := GetCacheService()
	ci := FindOrNewFuncCache(cs, ck)
	cs.AddOpenScope(ci)
	defer cs.CloseScope(ci)

	if ci.IsValid() {
		return GetVal[V](ci)
	}

	val, err := f(arg1, arg2)
	if err != nil {
		return defaultOut, err
	}

	ci.SetValue(val)
	return val, nil
}

func CacheCall3Args[
	a1 any,
	a2 any,
	a3 any,
	V any,
](f func(a1, a2, a3) (V, error), arg1 a1, arg2 a2, arg3 a3) (V, error) {
	var defaultOut V

	ck, err := newCacheKey(f, arg1, arg2, arg3)
	if err != nil {
		return defaultOut, err
	}

	cs := GetCacheService()
	ci := FindOrNewFuncCache(cs, ck)
	cs.AddOpenScope(ci)
	defer cs.CloseScope(ci)

	if ci.IsValid() {
		return GetVal[V](ci)
	}

	val, err := f(arg1, arg2, arg3)
	if err != nil {
		return defaultOut, err
	}

	ci.SetValue(val)
	return val, nil
}
