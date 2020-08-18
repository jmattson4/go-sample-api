package news

import (
	"github.com/jinzhu/gorm"
)

type NewsRepo struct {
	db *gorm.DB
}

func ConstructRepo(db *gorm.DB) NewsRepository {
	return &NewsRepo{
		db: db,
	}
}

func (repo *NewsRepo) RollBack() error {
	if err := repo.db.Rollback().Error; err != nil {
		return err
	}
	return nil
}

//Create Params: news *model.NewsData; Returns: error; Description: Creates Scraped News Data piece
func (repo *NewsRepo) Create(news *NewsData) error {
	err := repo.db.Create(news).Error
	if err != nil {
		return err
	}
	return nil
}

//Update Params: news *model.NewsData; Returns: error; Description: Updates News Data piece
func (repo *NewsRepo) Update(news *NewsData) error {
	err := repo.db.Model(news).Update(map[string]interface{}{
		"articleLink":   news.ArticleLink,
		"articleText":   news.ArticleText,
		"imageURL":      news.ImageURL,
		"paragraph":     news.Paragraph,
		"websiteName":   news.WebsiteName,
		"ArticleNumber": news.ArticleNumber,
	}).Error
	return err
}
func (repo *NewsRepo) Delete(news *NewsData) error {
	err := repo.db.Delete(news).Error
	return err
}

func (repo *NewsRepo) HardDelete(news *NewsData) error {
	err := repo.db.Unscoped().Delete(news).Error
	return err
}

//Get Params: news *model.NewsData; Returns: error; Description: Gets News Data piece
func (repo *NewsRepo) Get(news *NewsData) error {
	err := repo.db.Where("id = ?", news.ID).First(news).Error
	return err
}

//GetByArticleLink Params: news *model.NewsData; Returns: error; Description: Gets News Data piece by its article link
func (repo *NewsRepo) GetByArticleLink(news *NewsData) error {
	err := repo.db.Where("article_link = ?", news.ArticleLink).First(news).Error
	return err
}
func (repo *NewsRepo) GetNewsByWebNameAndID(news *NewsData) error {
	err := repo.db.Where("website_name = ? AND id = ?", news.WebsiteName, news.ID).First(news).Error
	return err
}
func (repo *NewsRepo) GetMultipleNews(start int, count int, news []*NewsData) error {
	err := repo.db.Offset(start).Limit(count).Find(news).Error

	if err != nil {
		return err
	}

	return nil
}
func (repo *NewsRepo) GetMultipleNewsByWebName(start, count int, news []*NewsData, webName string) error {
	err := repo.db.Offset(start).Limit(count).Where("website_name = ?", webName).Find(news).Error

	if err != nil {
		return err
	}

	return nil
}
