package scrapers

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//GlobalNewsData ... Struct is used to model the data taken from the GlobalNewsMainpage
type GlobalNewsData struct {
	ArticleLink []string
	ArticleText []string
	ImageURL    []string
	Paragraphs  []string
}

//Factory function for creation of GlobalNewsData Struct.
func newGlobalNewsData(articleLink []string, text []string, imageURL []string, paragraphs []string) *GlobalNewsData {

	data := GlobalNewsData{
		ArticleLink: articleLink,
		ArticleText: text,
		ImageURL:    imageURL,
		Paragraphs:  paragraphs,
	}

	return &data
}

//MainPageScrape ... This scrapes the mainpage of globalnew.ca returning
//	the links, text and img links of the top stories for the day.
func GlobalNewsScrape() *GlobalNewsData {
	doc, err := goquery.NewDocument("https://globalnews.ca")
	if err != nil {
		log.Fatal(err)
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

	linkSlice = unique(linkSlice)
	textSlice = unique(textSlice)
	imgSlice = unique(imgSlice)
	paragraphSlice = scrapeArticles(linkSlice)

	return newGlobalNewsData(linkSlice, textSlice, imgSlice, paragraphSlice)
}

func scrapeArticles(articleLinkSlice []string) []string {
	articleParagraphs := []string{}

	for i := 1; i < len(articleLinkSlice); i++ {
		paragraph := articleScrape(articleLinkSlice[i])
		articleParagraph := articleToParagraph(paragraph)
		articleParagraphs = append(articleParagraphs, articleParagraph)
	}
	return articleParagraphs
}

func articleScrape(linkSlice string) []string {
	doc, err := goquery.NewDocument(linkSlice)
	if err != nil {
		log.Fatal(err)
	}
	paragraphSlice := []string{}
	article := doc.Find("article")
	article.Find("p").Each(func(index int, item *goquery.Selection) {
		text := strings.TrimSpace(item.Text())
		paragraphSlice = append(paragraphSlice, text)
	})
	return paragraphSlice
}

func articleToParagraph(paragraphs []string) string {
	var paragraph string
	for _, pString := range paragraphs {
		paragraph = paragraph + " " + pString
	}
	return paragraph
}

//Function makes sure that everything that is pulled back is unique.
func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
