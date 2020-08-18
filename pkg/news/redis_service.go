package news

import (
	"errors"
)

type NewsCacheServ struct {
	Repo BaseRepository
}

func ConstructCacheService(r NewsRepository) *NewsServ {
	return &NewsServ{
		Repo: r,
	}
}

//Create Params: news *model.NewsData; Returns: error; Description: Creates Scraped News Data piece
func (serv *NewsCacheServ) Create(news *NewsData) error {
	err := serv.Repo.Create(news)
	if err != nil {
		return err
	}
	return nil
}

//Update Params: news *model.NewsData; Returns: error; Description: Updates News Data piece
func (serv *NewsCacheServ) Update(news *NewsData) error {
	vErr := validateID(news.ID)
	if vErr != nil {
		return vErr
	}
	err := serv.Repo.Update(news)
	if err != nil {
		return err
	}
	return nil
}
func (serv *NewsCacheServ) Delete(news *NewsData) error {
	vErr := validateID(news.ID)
	if vErr != nil {
		return vErr
	}
	err := serv.Repo.Delete(news)
	if err != nil {
		return err
	}
	return nil
}

//Get Params: news *model.NewsData; Returns: error; Description: Gets News Data piece
func (serv *NewsCacheServ) Get(news *NewsData) error {
	vErr := validateID(news.ID)
	if vErr != nil {
		return vErr
	}
	err := serv.Repo.Get(news)
	if err != nil {
		return err
	}
	return nil
}

func (serv *NewsCacheServ) GetMultipleNews(start int, count int, news []*NewsData) error {
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

	err := serv.Repo.GetMultipleNews(start, count, news)

	if err != nil {
		return err
	}

	return nil
}
