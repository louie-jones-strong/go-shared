package cache

var cacheServiceInstance *CacheService

type CacheInstance interface {
	IsValid() bool
	AddSubScope(sub CacheInstance)
	GetSubScopes() []CacheInstance
	ToString() string
	Clear()
}

type CacheService struct {
	openScopeStack []CacheInstance
	funcCaches     map[cacheKey]*FuncCache
	rootCaches     []CacheInstance
}

func GetCacheService() *CacheService {
	if cacheServiceInstance == nil {
		cacheServiceInstance = &CacheService{
			openScopeStack: make([]CacheInstance, 0, 10),
			funcCaches:     make(map[cacheKey]*FuncCache),
			rootCaches:     make([]CacheInstance, 0, 10),
		}
	}
	return cacheServiceInstance
}

func (s *CacheService) GetRootCaches() []CacheInstance {
	return s.rootCaches
}

func (s *CacheService) getLastOpenScope() CacheInstance {
	if len(s.openScopeStack) == 0 {
		return nil
	}
	return s.openScopeStack[len(s.openScopeStack)-1]
}

func (s *CacheService) AddOpenScope(scope CacheInstance) {
	last := s.getLastOpenScope()
	if last == nil {
		s.rootCaches = append(s.rootCaches, scope)
	} else {
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
