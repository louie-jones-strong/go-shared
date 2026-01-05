package cache

type FuncCache struct {
	val      any
	isValSet bool

	subInstances []CacheInstance
}

func NewFuncCache() *FuncCache {
	return &FuncCache{
		val:          nil,
		isValSet:     false,
		subInstances: nil,
	}
}

func (c *FuncCache) IsValid() bool {
	if !c.isValSet {
		return false
	}

	for _, instance := range c.subInstances {
		if !instance.IsValid() {
			return false
		}
	}
	return true
}

func (c *FuncCache) AddSubScope(sub CacheInstance) {
	c.subInstances = append(c.subInstances, sub)
}

func (c *FuncCache) GetValue() (any, bool) {
	if !c.IsValid() {
		return nil, false
	}
	return c.val, true
}

func (c *FuncCache) SetValue(value any) {
	c.val = value
	c.isValSet = true
}

func GetVal[V any](ci *FuncCache) (V, bool) {
	var defaultOut V

	val, hasVal := ci.GetValue()
	if !hasVal {
		return defaultOut, false
	}

	typeVal, ok := val.(V)
	if !ok {
		panic("cached value has unexpected type")
	}

	return typeVal, true
}
