package news

import (
	"github.com/jinzhu/gorm"
)

//NewsData .... Models Collected data from various scraped Canadian News Sites
type NewsData struct {
	gorm.Model
	ArticleLink   string `json:"articleLink" gorm:"not null"`
	ArticleText   string `json:"articleText"`
	ImageURL      string `json:"imageURL"`
	Paragraph     string `json:"paragraph"`
	WebsiteName   string `json:"websiteName"`
	ArticleNumber uint   `json:"articleNumber"`
}

//NewsDataBasicInit ...
//Description: Simple Factory
func NewsDataBasicInit() *NewsData {
	news := &NewsData{}
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
