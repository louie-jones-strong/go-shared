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
