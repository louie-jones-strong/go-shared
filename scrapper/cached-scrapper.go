package scrapper

import (
	"time"

	"github.com/louie-jones-strong/go-shared/filecache"
)

type CachedScrapper struct {
	siteCache *filecache.FileCache[string]
	scrapper  Scrapper
	fileExt   string
}

func NewCachedScrapper(
	siteCache *filecache.FileCache[string],
	fileExt string,
) *CachedScrapper {

	return &CachedScrapper{
		siteCache: siteCache,
		fileExt:   fileExt,
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
	fi := cs.siteCache.GetOrCreateFileInfo(url)
	if fi.IsValid(expireDuration) {
		return fi.LoadFile(filecache.DefaultFileKey)
	}

	// fall back to scrapping
	data, err := cs.ScrapURL(url)
	if err != nil {
		return nil, err
	}

	// save to cache
	err = fi.SaveFile(filecache.DefaultFileKey, cs.fileExt, data)
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
