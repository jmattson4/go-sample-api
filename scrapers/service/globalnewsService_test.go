package service_test

import (
	"testing"

	"github.com/jmattson4/go-sample-api/scrapers/service"
)

func TestGlobalNewsMainPageScrape(t *testing.T) {
	scrapeData, err := service.GlobalNewsScrape()
	if err != nil {
		t.Errorf("Scrape came back with an error: %v", err)
	}
	checkScrapeSize(t, 30, len(scrapeData.ArticleLink))
	checkScrapeSize(t, 30, len(scrapeData.ImageURL))
	checkScrapeSize(t, 30, len(scrapeData.Paragraphs))
	checkScrapeSize(t, 30, len(scrapeData.ArticleText))
}

func checkScrapeSize(t *testing.T, expected int, actual int) {
	if expected > actual {
		t.Errorf("Expected Scrape Size to be greater than %d. Got %d\n", expected, actual)
	}
}
