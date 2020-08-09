package model

import (
	"github.com/jinzhu/gorm"
)

//NewsData ....
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

//Create ...
//	Description: Creates the NewsData in the database given the current values of
//	the news data structure.
func (nd *NewsData) Create() error {
	err := GetDB().Create(nd).Error
	if err != nil {
		return err
	}
	return nil
}

// Get ...
// Description: Gets Newsdata given the newsdata ID
func (nd *NewsData) Get() error {
	err := GetDB().Where("id = ?", nd.ID).First(nd).Error
	return err
}

// GetNewsByArticleLink : This function grabs the First News with the associated article link
func (nd *NewsData) GetNewsByArticleLink() error {
	err := GetDB().Where("article_link = ?", nd.ArticleLink).First(nd).Error
	return err
}

// GetNewsByWebNameAndID : This function grabs the First News Article by website namd and article Link
func (nd *NewsData) GetNewsByWebNameAndID() error {
	err := GetDB().Where("website_name = ? AND id = ?", nd.WebsiteName, nd.ID).First(nd).Error
	return err
}

// GetMultipleNews ...
// Description: Given  a start, a count and a Newsdata Slice pointer it grabs the amount
//	within the start -> count interval.
func GetMultipleNews(start, count int, p *[]NewsData) error {

	err := GetDB().Offset(start).Limit(count).Find(p).Error

	if err != nil {
		return err
	}

	return nil
}

//GetMultipleNewsByWebName ...
//Description: Essentially get multiple but also passed a WebName to use in the query.
func GetMultipleNewsByWebName(start, count int, p *[]NewsData, webName string) error {

	err := GetDB().Offset(start).Limit(count).Where("website_name = ?", webName).Find(p).Error

	if err != nil {
		return err
	}

	return nil
}

// Update ...
// Description: This function does an update of NewsData that is grabbed
func (nd *NewsData) Update() error {
	err := GetDB().Model(nd).Update(map[string]interface{}{
		"articleLink": nd.ArticleLink,
		"articleText": nd.ArticleText,
		"imageURL":    nd.ImageURL,
		"paragraph":   nd.Paragraph,
	}).Error
	return err
}

// Delete ...
// Description: this function is uses the Gorm DB to soft delete
//	meaning it sets the "DeletedAt" field rather then actually delete
func (nd *NewsData) Delete() error {
	err := GetDB().Delete(nd).Error
	return err
}

// HardDelete ...
// Description: this function uses Gorm to actually do a hard delete.
//	Meaning it does a full delete from the database rather unlike Delete()
func (nd *NewsData) HardDelete() error {
	err := GetDB().Unscoped().Delete(nd).Error
	return err
}

//GetDeleted ...
//This function retrieves soft deleted items within the given start -> count interval.
//TODO : Make this grab the news_data_Delted index for optimization
func GetDeleted(start, count int, p *[]Product) error {
	err := GetDB().Unscoped().Offset(start).Limit(count).Find(p).Error

	if err != nil {
		return err
	}

	return nil
}

//ProcessNewsData : This function is used to process News Data from various sorces and save them into the Database.
func ProcessNewsData(
	processAmount uint,
	websiteName string,
	articleLink []string,
	articleText []string,
	imageURL []string,
	paragraphs []string) error {

	if processAmount > uint(len(articleLink)) {
		processAmount = uint(len(articleLink))
	}
	newsDataSlice := []*NewsData{}
	for i := uint(1); i < processAmount; i++ {
		newsData := NewsDataInit(articleLink[i], articleText[i], imageURL[i], paragraphs[i-1], websiteName, i)
		newsDataSlice = append(newsDataSlice, newsData)
	}
	for i := 0; i < len(newsDataSlice); i++ {
		newsD := newsDataSlice[i]
		err := GetDB().Create(newsD).Error
		if err != nil {
			GetDB().Rollback()
			return err
		}
	}
	return nil
}
