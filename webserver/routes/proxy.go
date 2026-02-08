package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/louie-jones-strong/go-shared/filecache"
	"github.com/louie-jones-strong/go-shared/logger"
)

func SetupProxyRouter(r *chi.Mux, proxyCache *filecache.FileCache[string]) {

	r.Get("/proxy/clear/*", func(w http.ResponseWriter, r *http.Request) {
		forwardURL, err := parseURL(chi.URLParam(r, "*"))
		if err != nil {
			logger.Error("Error parsing forward URL: %v", err)
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		numRemoved, err := proxyCache.RemoveFiles(
			forwardURL,
		)
		if err != nil {
			logger.Error("Error clearing cache for URL %s: %v", forwardURL, err)
			http.Error(w, "Error clearing cache", http.StatusInternalServerError)
			return
		}

		if numRemoved > 0 {
			logger.Info("Cleared %d cache items for URL: %s", numRemoved, forwardURL)
			w.WriteHeader(http.StatusOK)
			return
		}

		logger.Info("No cache items to clear for URL: %s", forwardURL)
		w.WriteHeader(http.StatusNotModified)
	})

	r.Get("/proxy/*", func(w http.ResponseWriter, r *http.Request) {
		forwardURL, err := parseURL(chi.URLParam(r, "*"))
		if err != nil {
			logger.Error("Error parsing forward URL: %v", err)
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		statusCode := http.StatusOK
		lastModified := time.Now().UTC()
		var (
			header http.Header
			body   []byte
		)

		fgi := proxyCache.TryGetFileInfo(forwardURL)

		if fgi != nil {

			lastModified = fgi.GetLastUpdated()

			clientModStr := r.Header.Get("If-Modified-Since")
			if clientModStr != "" {
				clientModTime, err := time.Parse(http.TimeFormat, clientModStr)
				if err == nil && !lastModified.After(clientModTime) {
					// Client cache is valid so we can return 304 Not Modified
					logger.Debug("Cache HIT Not Modified for URL: %s", forwardURL)
					w.WriteHeader(http.StatusNotModified)
					w.Header().Set("Cache-Control", "public")
					w.Header().Set("Last-Modified", lastModified.Format(http.TimeFormat))
					return
				}
			}

			// need to send sever cache to client
			logger.Debug("Cache HIT for URL: %s", forwardURL)
			header, body, err = tryLoadResponseFromCache(fgi)
			if err != nil {
				http.Error(w, "Cache read error", http.StatusInternalServerError)
				return
			}

		} else {
			// not in server cache, do request
			logger.Debug("Cache MISS for URL: %s", forwardURL)
			statusCode, header, body, err = doRequest(forwardURL, r, proxyCache)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		removeCacheHeaderEntries(header, w.Header())
		w.Header().Set("Cache-Control", "public")
		w.Header().Set("Last-Modified", lastModified.Format(http.TimeFormat))

		w.WriteHeader(statusCode)
		_, err = w.Write(body)
		if err != nil {
			logger.Error("Error writing response body to client: %v", err)
			return
		}
	})
}

func tryLoadResponseFromCache(
	fgi *filecache.FileGroupInfo,
) (http.Header, []byte, error) {

	// load body from cache
	body, err := fgi.LoadFile("body")
	if err != nil {
		return nil, nil, fmt.Errorf("Error loading body from cache: %w", err)
	}

	// load header from cache
	headerData, err := fgi.LoadFile("header")
	if err != nil {
		return nil, nil, fmt.Errorf("Error loading header from cache: %w", err)
	}

	header := http.Header{}
	if headerData != nil {
		err = json.Unmarshal(headerData, &header)
		if err != nil {
			return nil, nil, fmt.Errorf("Error unmarshaling header JSON from cache: %w", err)
		}
	}

	return header, body, nil
}

func doRequest(parsedURL string, r *http.Request, cache *filecache.FileCache[string]) (int, http.Header, []byte, error) {

	requestStartTime := time.Now()

	req, err := http.NewRequest(r.Method, parsedURL, r.Body)
	if err != nil {
		return 500, nil, nil, fmt.Errorf("Failed to create request")
	}
	removeCacheHeaderEntries(r.Header, req.Header)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return 500, nil, nil, fmt.Errorf("Failed: %w", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 500, nil, nil, fmt.Errorf("Failed to read response body: %w", err)
	}

	requestDuration := time.Since(requestStartTime)
	logger.Info("Request to %s took %v", parsedURL, requestDuration)

	if res.StatusCode < http.StatusMultipleChoices {
		// Save to cache asynchronously to avoid blocking the response
		go func(cacheKey string, headerCopy http.Header, bodyCopy []byte) {
			saveStartTime := time.Now()
			saveErr := saveResponseToCache(cache, cacheKey, headerCopy, bodyCopy)
			if saveErr != nil {
				logger.Error("Error saving to cache for key %s: %v", cacheKey, saveErr)
			}
			saveDuration := time.Since(saveStartTime)
			logger.Debug("Saved response to cache for key %s in %v", cacheKey, saveDuration)
		}(parsedURL, res.Header.Clone(), body)
	}

	return res.StatusCode, res.Header, body, nil
}

func removeCacheHeaderEntries(h http.Header, res http.Header) {
	cacheKeysToSkip := []string{
		"cache-control",
		"last-modified",
		"x-cache-status",
		"age",
		"x-cache",
		"x-content-type-options",
		"x-served-by",
		"etag",
	}

	for key, values := range h {
		skip := false
		for _, skipKey := range cacheKeysToSkip {
			if strings.ToLower(key) == skipKey {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		for _, value := range values {
			res.Add(key, value)
		}
	}
}

func saveResponseToCache(
	cache *filecache.FileCache[string],
	urlStr string,
	header http.Header,
	body []byte,
) error {

	fgi := cache.GetOrCreateFileInfo(urlStr)

	// save body to cache
	ext := path.Ext(urlStr)
	err := fgi.SaveFile("body", ext, body)
	if err != nil {
		return fmt.Errorf("Error saving body to cache for key %s: %w", urlStr, err)
	}

	// save headers to cache
	// serialize headers to json then to []byte
	headerData, err := json.Marshal(header)
	if err != nil {
		return fmt.Errorf("Error marshaling headers to JSON for key %s: %w", urlStr, err)
	}

	err = fgi.SaveFile("header", "", headerData)
	if err != nil {
		return fmt.Errorf("Error saving headers to cache for key %s: %w", urlStr, err)
	}

	return nil
}

func parseURL(rawURL string) (string, error) {
	if rawURL == "" {
		return "", fmt.Errorf("parseURL called with empty URL")
	}

	rawURL = strings.TrimSpace(rawURL)
	rawURL = strings.ReplaceAll(rawURL, "% ", "%25"+"%20")
	rawURL = strings.ReplaceAll(rawURL, " ", "%20")

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	parsedURLStr := parsedURL.String()
	if parsedURLStr == "" {
		return "", fmt.Errorf("invalid URL: %s", rawURL)
	}
	return parsedURLStr, nil
}
