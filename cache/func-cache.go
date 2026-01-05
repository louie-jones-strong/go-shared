package cache

type FuncCache struct {
	BaseCacheScope[any]
	key cacheKey
}

func FindOrNewFuncCache(cs *CacheService, ck cacheKey) *FuncCache {
	cacheInstance, exists := cs.funcCaches[ck]
	if exists {
		return cacheInstance
	}

	c := &FuncCache{
		BaseCacheScope: NewBaseCacheScope[any](),
		key:            ck,
	}
	cs.funcCaches[ck] = c
	return c
}

func (c *FuncCache) ToString() string {
	return "FuncCache: " + c.key.ToString()
}

func GetVal[V any](ci *FuncCache) (V, bool) {

	val := ci.GetValue()
	typeVal, ok := val.(V)
	if !ok {
		panic("cached value has unexpected type")
	}

	return typeVal, true
}
