package scrapers_test

import (
	"testing"

	"github.com/jmattson4/go-sample-api/scrapers"
)

func TestMainPageScrape(t *testing.T) {
	scrapeData := scrapers.MainPageScrape()
	checkScrapeSize(t, 50, len(scrapeData.ArticleLink))
	checkScrapeSize(t, 50, len(scrapeData.ImageURL))
	checkScrapeSize(t, 50, len(scrapeData.Paragraphs))
	checkScrapeSize(t, 50, len(scrapeData.ArticleText))
}

func checkScrapeSize(t *testing.T, expected int, actual int) {
	if expected > actual {
		t.Errorf("Expected Scrape Size to be greater than %d. Got %d\n", expected, actual)
	}
}
