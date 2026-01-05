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

	val, hasVal := GetVal[V](ci)
	if hasVal {
		return val, nil
	}

	val, err = f()
	if err != nil {
		return defaultOut, err
	}

	ci.SetValue(val)
	return val, nil
}
