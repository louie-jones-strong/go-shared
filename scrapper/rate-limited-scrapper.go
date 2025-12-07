package scrapper

import (
	"net/url"
	"time"

	"github.com/louie-jones-strong/go-shared/logger"
)

type RateLimitedScrapper struct {
	minTimeBetweenRequests time.Duration
	scrapper               Scrapper
	hostLastScrappedTime   map[string]time.Time
}

func NewRateLimitedScrapper(
	minTimeBetweenRequests time.Duration,
	scrapper Scrapper,
) *RateLimitedScrapper {
	return &RateLimitedScrapper{
		minTimeBetweenRequests: minTimeBetweenRequests,
		scrapper:               scrapper,
		hostLastScrappedTime:   map[string]time.Time{},
	}
}

func (s *RateLimitedScrapper) CleanUp() {
	s.scrapper.CleanUp()
}

func (s *RateLimitedScrapper) ScrapURL(urlToScrap string) ([]byte, error) {
	parsedURL, err := url.Parse(urlToScrap)
	if err != nil {
		return nil, err
	}
	hostName := parsedURL.Hostname()

	lastScrappedTime, found := s.hostLastScrappedTime[hostName]
	if found {
		now := time.Now().UTC()
		timeSinceLast := now.Sub(lastScrappedTime)
		if timeSinceLast < s.minTimeBetweenRequests {
			waitDuration := s.minTimeBetweenRequests - timeSinceLast

			logger.Debug("Delaying request by %v", waitDuration)
			time.Sleep(waitDuration)
		}
	}
	s.hostLastScrappedTime[hostName] = time.Now().UTC()

	return s.scrapper.ScrapURL(urlToScrap)
}
