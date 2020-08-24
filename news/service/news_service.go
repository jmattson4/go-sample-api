package service

import (
	"errors"
	"fmt"

	"github.com/twinj/uuid"

	"github.com/jmattson4/go-sample-api/domain"
)

type NewsServ struct {
	db    domain.NewsDBRepository
	cache domain.NewsCacheRepository
}

func ConstructService(r domain.NewsDBRepository, c domain.NewsCacheRepository) *NewsServ {
	return &NewsServ{
		db:    r,
		cache: c,
	}
}

//Create Params: news *model.NewsData; Returns: error; Description: Creates Scraped News Data piece
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

//Update Params: news *model.NewsData; Returns: error; Description: Updates News Data piece
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
func (serv *NewsServ) GetMultipleNews(start int, count int, news []*domain.NewsData) error {
	if start < 0 {
		err := errors.New("Start must be greater than zero.")
		return err
	}
	if count < 0 {
		err := errors.New("Count must be greater than zero.")
		return err
	}
	if len(news) < 0 {
		err := errors.New("Length of News Slice must be greater than 0")
		return err
	}

	cacheErr := serv.cache.GetMultipleNews(start, count, news)
	if cacheErr != nil {
		err := serv.db.GetMultipleNews(start, count, news)
		if err != nil {
			return err
		}
	}
	return nil
}
func (serv *NewsServ) GetMultipleNewsByWebName(start, count int, news []*domain.NewsData, webName string) error {
	if start < 0 {
		err := errors.New("Start must be greater than zero.")
		return err
	}
	if count < 0 {
		err := errors.New("count must be greater than zero.")
		return err
	}

	err := serv.db.GetMultipleNewsByWebName(start, count, news, webName)

	if err != nil {
		return err
	}

	return nil
}

func (s *NewsServ) ProcessNewsData(
	processAmount uint,
	websiteName string,
	articleLink []string,
	articleText []string,
	imageURL []string,
	paragraphs []string) error {

	if processAmount > uint(len(articleLink)) {
		processAmount = uint(len(articleLink))
	}
	newsDataSlice := []*domain.NewsData{}
	for i := uint(1); i < processAmount; i++ {
		newsData := domain.NewsDataInit(articleLink[i], articleText[i], imageURL[i], paragraphs[i-1], websiteName, i)
		newsDataSlice = append(newsDataSlice, newsData)
	}
	for i := 0; i < len(newsDataSlice); i++ {
		newsD := newsDataSlice[i]
		err := s.db.Create(newsD)
		if err != nil {
			s.db.RollBack()
			return err
		}
	}
	return nil

}

//Validation

func validateID(id uuid.UUID) error {
	if uuid.IsNil(id) == true {
		err := errors.New("News ID cannot be null or empty")
		return err
	}
	return nil
}
func validateString(s string, propertyName string) error {
	if len(s) <= 0 {
		err := errors.New(fmt.Sprintf("%v cannot be empty.", propertyName))
		return err
	}
	return nil
}
