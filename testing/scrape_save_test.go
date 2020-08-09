package testing

import (
	"testing"

	"github.com/jmattson4/go-sample-api/model"
	"github.com/jmattson4/go-sample-api/scrapers"
)

func TestScrapeAndSave(t *testing.T) {
	scrapeData := scrapers.GlobalNewsScrape()
	err := model.ProcessNewsData(20, "globalnews", scrapeData.ArticleLink, scrapeData.ArticleText, scrapeData.ImageURL, scrapeData.Paragraphs)
	if err != nil {
		t.Errorf("Failed To Scrape and Save Test Data! Reason: %v", err)
	}
	scrapeTestFind := &model.NewsData{}
	scrapeTestFind.ArticleLink = scrapeData.ArticleLink[1]
	getErr := scrapeTestFind.GetNewsByArticleLink()
	if getErr != nil {
		t.Errorf("Failed To Get Newly created news By Article Link! Reason: %v", getErr)
	}

}
