package cache

type BaseCacheScope[T any] struct {
	subInstances []CacheInstance
	cachedVal    *T
}

func NewBaseCacheScope[T any]() BaseCacheScope[T] {
	return BaseCacheScope[T]{
		subInstances: make([]CacheInstance, 0, 10),
		cachedVal:    new(T),
	}
}

func (c *BaseCacheScope[T]) IsValid() bool {
	if c.cachedVal == nil {
		return false
	}

	isValid := true
	for _, instance := range c.subInstances {
		if !instance.IsValid() {
			isValid = false
			break
		}
	}

	if !isValid {
		c.Clear()
	}
	return isValid
}

func (c *BaseCacheScope[T]) AddSubScope(sub CacheInstance) {
	c.subInstances = append(c.subInstances, sub)
}

func (c *BaseCacheScope[T]) GetSubScopes() []CacheInstance {
	return c.subInstances
}

func (c *BaseCacheScope[T]) GetValue() T {
	return *c.cachedVal
}

func (c *BaseCacheScope[T]) SetValue(value T) {
	c.cachedVal = &value
}

func (c *BaseCacheScope[T]) Clear() {
	c.cachedVal = nil
}
