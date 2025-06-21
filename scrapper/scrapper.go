package scrapper

type Scrapper interface {
	CleanUp()
	ScrapURL(url string) ([]byte, error)
}
