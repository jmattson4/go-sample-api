package domain

type RawNewsData struct {
	ArticleLink []string
	ArticleText []string
	ImageURL    []string
	Paragraphs  []string
}

func NewRawNewsData(articleLink []string, text []string, imageURL []string, paragraphs []string) *RawNewsData {
	return &RawNewsData{
		ArticleLink: articleLink,
		ArticleText: text,
		ImageURL:    imageURL,
		Paragraphs:  paragraphs,
	}
}

type RawNewsService interface {
	RawNewsScrape() (*RawNewsData, error)
}
