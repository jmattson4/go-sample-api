package domain

//NewsData .... Models Collected data from various scraped Canadian News Sites
type NewsData struct {
	Base
	ArticleLink   string `json:"articleLink" gorm:"not null"`
	ArticleText   string `json:"articleText"`
	ImageURL      string `json:"imageURL"`
	Paragraph     string `json:"paragraph"`
	WebsiteName   string `json:"websiteName"`
	ArticleNumber uint   `json:"articleNumber"`
}

//NewsDataNoInit ...
func NewsDataNoInit() *NewsData {
	news := &NewsData{}
	return news
}

//NewsDataBasicInit ...
//Description: Simple Factory
func NewsDataBasicInit() *NewsData {
	base := Base{}
	news := &NewsData{
		Base: base,
	}
	return news
}

//NewsDataInit ...
///Description: factory function used to init  a NewsData struct
func NewsDataInit(al string, at string, img string, p string, webN string, artN uint) *NewsData {
	news := &NewsData{
		ArticleLink:   al,
		ArticleText:   at,
		ImageURL:      img,
		Paragraph:     p,
		WebsiteName:   webN,
		ArticleNumber: artN,
	}
	return news
}

/*
	Below is the interface declartions for the news section of the application.
	It contains definitions for base Reading and Writing to external Datasources
	It also contains some other classes used for The News repository and Service that
	Connects to PostGreSQL. Can be expanded to fit other External Datasources.

*/

type NewsCacheRepository interface {
	Get(news *NewsData) error
	Create(news *NewsData) error
	Update(news *NewsData) error
	Delete(news *NewsData) error
}

type NewsDBRepository interface {
	Get(news *NewsData) error
	GetMultipleNews(start int, count int, news []*NewsData) error
	Create(news *NewsData) error
	Update(news *NewsData) error
	Delete(news *NewsData) error
	GetMultipleNewsByWebName(start, count int, news []*NewsData, webName string) error
	HardDelete(news *NewsData) error
	GetByArticleLink(news *NewsData) error
	GetNewsByWebNameAndID(news *NewsData) error
	RollBack() error
}

type NewsService interface {
	Get(news *NewsData) error
	GetMultipleNews(start int, count int, news []*NewsData) error
	Create(news *NewsData) error
	Update(news *NewsData) error
	Delete(news *NewsData) error
	GetMultipleNewsByWebName(start, count int, news []*NewsData, webName string) error
	HardDelete(news *NewsData) error
	GetByArticleLink(news *NewsData) error
	GetNewsByWebNameAndID(news *NewsData) error
	ProcessNewsData(
		processAmount uint,
		websiteName string,
		raw *RawNewsData) error
}
