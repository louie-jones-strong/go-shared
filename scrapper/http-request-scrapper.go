package scrapper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type HTTPRequestScrapper struct {
}

func NewHTTPRequestScrapper() *HTTPRequestScrapper {

	s := &HTTPRequestScrapper{}
	return s
}

func (s *HTTPRequestScrapper) CleanUp() {
}

func (s *HTTPRequestScrapper) ScrapURL(url string) ([]byte, error) {
	start := time.Now()

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		log.Printf("failed response: %v", res)
		return nil, fmt.Errorf("Scraping URL: %v returned code: %v", url, res.StatusCode)
	}

	log.Printf("Scraping URL: %v Took: %v", url, time.Since(start))

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return bodyBytes, nil
}
