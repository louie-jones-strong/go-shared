package cache

func CacheCall[V any](f func() (V, error)) (V, error) {
	var defaultOut V

	ck, err := newCacheKey(f)
	if err != nil {
		return defaultOut, err
	}

	cs := GetCacheService()
	ci := cs.GetOrCreate(ck)

	val, hasVal := GetVal[V](ci)
	if hasVal {
		return val, nil
	}

	cs.AddOpenScope(ci)

	val, err = f()
	if err != nil {
		cs.CloseScope(ci)
		return defaultOut, err
	}

	ci.SetValue(val)
	cs.CloseScope(ci)
	return val, nil
}
