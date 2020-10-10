package service

import (
	"log"
	"strings"

	"github.com/jmattson4/go-sample-api/domain"

	"github.com/PuerkitoBio/goquery"
	u "github.com/jmattson4/go-sample-api/util"
)

type GlobalNewsScraperService struct {
}

//GlobalNewsScrape ...
//	Implementation of RawNewsService
//	This scrapes the mainpage of globalnew.ca returning
//	the links, text and img links of the top stories for the day.
func GlobalNewsScrape() (*domain.RawNewsData, error) {
	doc, err := goquery.NewDocument("https://globalnews.ca")
	if err != nil {
		log.Print(domain.SCRAPER_GLOBAL_NEWS_CANT_FIND_URL)
		return nil, domain.SCRAPER_GLOBAL_NEWS_CANT_FIND_URL
	}
	var textSlice []string
	var linkSlice []string
	var imgSlice []string
	var paragraphSlice []string

	doc.Find("div .c-posts__inner").Each(func(index int, item *goquery.Selection) {
		anchorTag := item.Find("a").Each(func(index int, item *goquery.Selection) {
			imgTag := item.Find("img")
			imgURL, _ := imgTag.Attr("data-src")
			imgSlice = append(imgSlice, imgURL)
		})

		linkHref, _ := anchorTag.Attr("href")
		linkText := strings.TrimSpace(anchorTag.Text())

		textSlice = append(textSlice, linkText)
		linkSlice = append(linkSlice, linkHref)
	})

	linkSlice = u.Unique(linkSlice)
	textSlice = u.Unique(textSlice)
	imgSlice = u.Unique(imgSlice)
	paragraphSlice, err = scrapeArticles(linkSlice)

	if err != nil {
		return nil, err
	}

	return domain.NewRawNewsData(linkSlice, textSlice, imgSlice, paragraphSlice), nil
}

//takes all the article links from the initial page scrape then goes to the pages scraping more data relevant to the article,
func scrapeArticles(articleLinkSlice []string) ([]string, error) {
	articleParagraphs := []string{}

	for i := 1; i < len(articleLinkSlice); i++ {
		paragraph, err := articleScrape(articleLinkSlice[i])
		if err != nil {
			return nil, err
		}
		articleParagraph := articleToParagraph(paragraph)
		articleParagraphs = append(articleParagraphs, articleParagraph)
	}
	return articleParagraphs, nil
}

func articleScrape(linkSlice string) ([]string, error) {
	doc, err := goquery.NewDocument(linkSlice)
	if err != nil {
		log.Print(domain.SCRAPER_GLOBAL_NEWS_CANT_FIND_URL.Error())
		return nil, domain.SCRAPER_GLOBAL_NEWS_CANT_FIND_URL
	}
	paragraphSlice := []string{}
	article := doc.Find("article")
	article.Find("p").Each(func(index int, item *goquery.Selection) {
		text := strings.TrimSpace(item.Text())
		paragraphSlice = append(paragraphSlice, text)
	})
	return paragraphSlice, nil
}

func articleToParagraph(paragraphs []string) string {
	var paragraph string
	for _, pString := range paragraphs {
		paragraph = paragraph + " " + pString
	}
	return paragraph
}
