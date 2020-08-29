package domain

import (
	"github.com/jinzhu/gorm"
	"github.com/twinj/uuid"
)

//NewsData .... Models Collected data from various scraped Canadian News Sites
type NewsData struct {
	gorm.Model
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ArticleLink   string    `json:"articleLink" gorm:"not null"`
	ArticleText   string    `json:"articleText"`
	ImageURL      string    `json:"imageURL"`
	Paragraph     string    `json:"paragraph"`
	WebsiteName   string    `json:"websiteName"`
	ArticleNumber uint      `json:"articleNumber"`
}

func NewsDataNoInit() *NewsData {
	news := &NewsData{}
	return news
}

//NewsDataBasicInit ...
//Description: Simple Factory
func NewsDataBasicInit() *NewsData {
	uuid := uuid.NewV4()
	news := &NewsData{
		ID: uuid,
	}
	return news
}

//NewsDataInit ...
///Description: factory function used to init  a NewsData struct
func NewsDataInit(al string, at string, img string, p string, webN string, artN uint) *NewsData {
	uuid := uuid.NewV4()
	news := &NewsData{
		ID:            uuid,
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
