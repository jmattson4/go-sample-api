package service

import (
	"fmt"

	"github.com/twinj/uuid"

	"github.com/jmattson4/go-sample-api/domain"
)

//NewsServ ...
// Description: struct which holds references to a db repo and a cache repo
type NewsServ struct {
	db    domain.NewsDBRepository
	cache domain.NewsCacheRepository
}

//ConstructNewsService ...
//	Params: r domain.NewsDBRepository, c domain.NewsCacheRepository.
//	Returns: *NewsServ
// 	Description: Factory function to create a NewsServ struct.
func ConstructNewsService(r domain.NewsDBRepository, c domain.NewsCacheRepository) *NewsServ {
	return &NewsServ{
		db:    r,
		cache: c,
	}
}

//Create ...
//	Params: news *model.NewsData;
//	Returns: error;
//	Description: Creates a row in the database for a given newsData piece
func (serv *NewsServ) Create(news *domain.NewsData) error {
	err := serv.db.Create(news)
	if err != nil {
		return err
	}
	cacheErr := serv.cache.Create(news)
	if cacheErr != nil {
		return cacheErr
	}
	return nil
}

//Update ...
//	Params: news *model.NewsData;
//	Returns: error;
//	Description: Updates News Data piece
func (serv *NewsServ) Update(news *domain.NewsData) error {
	vErr := validateID(news.ID)
	if vErr != nil {
		return vErr
	}
	err := serv.db.Update(news)
	if err != nil {
		return err
	}
	cacheErr := serv.cache.Update(news)
	if cacheErr != nil {
		return cacheErr
	}
	return nil
}

//Delete ...
// This deletes the news item from the cache and the DB.
func (serv *NewsServ) Delete(news *domain.NewsData) error {
	vErr := validateID(news.ID)
	if vErr != nil {
		return vErr
	}
	err := serv.db.Delete(news)
	if err != nil {
		return err
	}
	cacheErr := serv.cache.Delete(news)
	if cacheErr != nil {
		return cacheErr
	}
	return nil
}

//HardDelete ...
//	Pretty much the same implementation of Delete. Kind of need to possibly get rid of this function.
func (serv *NewsServ) HardDelete(news *domain.NewsData) error {
	vErr := validateID(news.ID)
	if vErr != nil {
		return vErr
	}
	err := serv.db.HardDelete(news)
	if err != nil {
		return err
	}
	cacheErr := serv.cache.Delete(news)
	if cacheErr != nil {
		return cacheErr
	}
	return nil
}

//Get Params: news *model.NewsData; Returns: error; Description: Gets News Data piece
func (serv *NewsServ) Get(news *domain.NewsData) error {
	vErr := validateID(news.ID)
	if vErr != nil {
		return vErr
	}
	err := serv.cache.Get(news)
	if err != nil {
		dbErr := serv.db.Get(news)
		if dbErr != nil {
			return dbErr
		}
		serv.cache.Create(news)
	}
	return nil
}

//GetByArticleLink Params: news *model.NewsData; Returns: error; Description: Gets News Data piece by its article link
func (serv *NewsServ) GetByArticleLink(news *domain.NewsData) error {
	vErr := validateID(news.ID)
	if vErr != nil {
		return vErr
	}
	alErr := validateString(news.ArticleLink, "Article Link")
	if alErr != nil {
		return alErr
	}
	err := serv.db.GetByArticleLink(news)
	if err != nil {
		return err
	}
	return nil
}

//GetNewsByWebNameAndID ...
//	Params: news *domain.NewsData
//	Return: error
// Description:
func (serv *NewsServ) GetNewsByWebNameAndID(news *domain.NewsData) error {
	vErr := validateID(news.ID)
	if vErr != nil {
		return vErr
	}
	wnErr := validateString(news.WebsiteName, "Website Name")
	if wnErr != nil {
		return wnErr
	}
	err := serv.db.GetNewsByWebNameAndID(news)
	if err != nil {
		return err
	}
	return nil
}

//GetMultipleNews ...
//	Description: This gets multiple records from the database
func (serv *NewsServ) GetMultipleNews(start int, count int, news []*domain.NewsData) error {
	if start < 0 {
		err := domain.NEWS_SERVICE_GETMULTIPLENEWS_START
		return err
	}
	if count < 0 {
		err := domain.NEWS_SERVICE_GETMULTIPLENEWS_COUNT
		return err
	}
	if len(news) < 0 {
		err := domain.NEWS_SERVICE_GETMULTIPLENEWS_EMPTYSLICE
		return err
	}

	err := serv.db.GetMultipleNews(start, count, news)
	if err != nil {
		return err
	}
	return nil
}

//GetMultipleNewsByWebName ...
//	Description:
func (serv *NewsServ) GetMultipleNewsByWebName(start, count int, news []*domain.NewsData, webName string) error {
	if start < 0 {
		err := domain.NEWS_SERVICE_GETMULTIPLENEWS_START
		return err
	}
	if count < 0 {
		err := domain.NEWS_SERVICE_GETMULTIPLENEWS_COUNT
		return err
	}

	err := serv.db.GetMultipleNewsByWebName(start, count, news, webName)

	if err != nil {
		return err
	}

	return nil
}

//ProcessNewsData ...
// Description: This function will process a bulk amount of newsData and enter it into the
// database
func (serv *NewsServ) ProcessNewsData(
	processAmount uint,
	websiteName string,
	raw *domain.RawNewsData) error {

	if processAmount > uint(len(raw.ArticleLink)) {
		processAmount = uint(len(raw.ArticleLink))
	}
	newsDataSlice := []*domain.NewsData{}
	for i := uint(1); i < processAmount; i++ {
		newsData := domain.NewsDataInit(raw.ArticleLink[i], raw.ArticleText[i], raw.ImageURL[i], raw.Paragraphs[i-1], websiteName, i)
		newsDataSlice = append(newsDataSlice, newsData)
	}
	for i := 0; i < len(newsDataSlice); i++ {
		newsD := newsDataSlice[i]
		err := serv.db.Create(newsD)
		if err != nil {
			serv.db.RollBack()
			return err
		}
	}
	return nil

}

//Validation
//validateID ...
// Description: This function validates a uuid checking to see if it is nil
func validateID(id uuid.UUID) error {
	if uuid.IsNil(id) == true {
		err := domain.NEWS_SERVICE_VALIDATEID_UUID_ISNIL
		return err
	}
	return nil
}

//validateString ...
// Description: This function checks to see if a string is empty.
func validateString(s string, propertyName string) error {
	if len(s) <= 0 {
		err := fmt.Errorf("%v cannot be empty", propertyName)
		return err
	}
	return nil
}
