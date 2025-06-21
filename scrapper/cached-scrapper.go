package scrapper

import (
	"time"

	"github.com/louie-jones-strong/go-shared/filecache"
)

type CachedScrapper struct {
	siteCache *filecache.FileCache[string]
	scrapper  Scrapper
}

func NewCachedScrapper(
	siteCache *filecache.FileCache[string],
) *CachedScrapper {

	return &CachedScrapper{
		siteCache: siteCache,
	}
}

func (cs *CachedScrapper) SetScrapper(
	scrapper Scrapper,
) {
	cs.scrapper = scrapper
}

func (cs *CachedScrapper) CleanUp() {
	if cs.scrapper != nil {
		cs.scrapper.CleanUp()
		cs.scrapper = nil
	}
}

func (cs *CachedScrapper) ScrapURLWithCache(url string, expireDuration time.Duration) ([]byte, error) {

	// try load from cache
	cachedData, err := cs.siteCache.TryLoadFileWithExpire(url, expireDuration)
	if err != nil {
		return nil, err
	}

	if cachedData != nil {
		return cachedData, nil
	}

	// fall back to scrapping
	data, err := cs.ScrapURL(url)
	if err != nil {
		return nil, err
	}

	// save to cache
	err = cs.siteCache.SaveFileWithExt(url, data, ".html")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (cs *CachedScrapper) ScrapURL(url string) ([]byte, error) {

	scrapper := cs.GetOrCreateScrapper()

	htmlSource, err := scrapper.ScrapURL(url)
	if err != nil {
		return nil, err
	}

	return []byte(htmlSource), nil
}

func (cs *CachedScrapper) GetOrCreateScrapper() Scrapper {

	if cs.scrapper == nil {
		scrapper := New()
		cs.scrapper = scrapper
	}

	return cs.scrapper
}
