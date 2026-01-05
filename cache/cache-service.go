package cache

var cacheServiceInstance *CacheService

type CacheInstance interface {
	IsValid() bool
	AddSubScope(sub CacheInstance)
}

type CacheService struct {
	openScopeStack []CacheInstance
	funcCaches     map[cacheKey]*FuncCache
}

func GetCacheService() *CacheService {
	if cacheServiceInstance == nil {
		cacheServiceInstance = &CacheService{
			openScopeStack: make([]CacheInstance, 0, 10),
			funcCaches:     make(map[cacheKey]*FuncCache),
		}
	}
	return cacheServiceInstance
}

func (s *CacheService) GetOrCreate(ck cacheKey) *FuncCache {

	cacheInstance, exists := s.funcCaches[ck]
	if !exists {
		cacheInstance = NewFuncCache()
		s.funcCaches[ck] = cacheInstance
	}

	return cacheInstance
}

func (s *CacheService) getLastOpenScope() CacheInstance {
	if len(s.openScopeStack) == 0 {
		return nil
	}
	return s.openScopeStack[len(s.openScopeStack)-1]
}

func (s *CacheService) AddOpenScope(scope CacheInstance) {
	last := s.getLastOpenScope()
	if last != nil {
		last.AddSubScope(scope)
	}
	s.openScopeStack = append(s.openScopeStack, scope)
}

func (s *CacheService) CloseScope(scope CacheInstance) {

	last := s.getLastOpenScope()
	if last != scope {
		panic("Can only close the last opened scope")
	}
	s.openScopeStack = s.openScopeStack[:len(s.openScopeStack)-1]
}
