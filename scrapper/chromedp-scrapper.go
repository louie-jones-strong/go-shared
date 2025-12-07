package scrapper

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/louie-jones-strong/go-shared/logger"
)

type ChromedpScrapper struct {
	scrapperCtx  context.Context
	cleanUpFuncs []context.CancelFunc
}

func New() *ChromedpScrapper {
	cancelFuncs := []context.CancelFunc{}

	ctx, cancel := chromedp.NewExecAllocator(
		context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),                                 // Run headless
			chromedp.Flag("disable-blink-features", "AutomationControlled"), // Bypass bot detection
			chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
			chromedp.WindowSize(1920, 1080),
		)...,
	)
	cancelFuncs = append(cancelFuncs, cancel)

	ctx, cancel = chromedp.NewContext(ctx)
	cancelFuncs = append(cancelFuncs, cancel)

	// ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	// cancelFuncs = append(cancelFuncs, cancel)

	s := &ChromedpScrapper{
		scrapperCtx:  ctx,
		cleanUpFuncs: cancelFuncs,
	}
	return s
}

func (s *ChromedpScrapper) CleanUp() {
	for _, cleanUpFunc := range s.cleanUpFuncs {
		cleanUpFunc()
	}
}

func (s *ChromedpScrapper) ScrapURL(url string) ([]byte, error) {
	start := time.Now()

	var htmlSource string
	err := s.RunActions(
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &htmlSource),
	)
	if err != nil {
		return nil, err
	}

	logger.Debug("Scraping URL: %v Took: %v", url, time.Since(start))

	return []byte(htmlSource), nil
}

func (s *ChromedpScrapper) RunActions(actions ...chromedp.Action) error {

	err := chromedp.Run(s.scrapperCtx,
		actions...,
	)
	if err != nil {
		return err
	}

	return nil
}
