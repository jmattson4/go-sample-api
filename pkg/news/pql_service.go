package news

import (
	"errors"
)

type NewsServ struct {
	Repo NewsRepository
}

func ConstructService(r NewsRepository) *NewsServ {
	return &NewsServ{
		Repo: r,
	}
}

//Create Params: news *model.NewsData; Returns: error; Description: Creates Scraped News Data piece
func (serv *NewsServ) Create(news *NewsData) error {
	err := serv.Repo.Create(news)
	if err != nil {
		return err
	}
	return nil
}

//Update Params: news *model.NewsData; Returns: error; Description: Updates News Data piece
func (serv *NewsServ) Update(news *NewsData) error {
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
func (serv *NewsServ) Delete(news *NewsData) error {
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

func (serv *NewsServ) HardDelete(news *NewsData) error {
	vErr := validateID(news.ID)
	if vErr != nil {
		return vErr
	}
	err := serv.Repo.HardDelete(news)
	if err != nil {
		return err
	}
	return nil
}

//Get Params: news *model.NewsData; Returns: error; Description: Gets News Data piece
func (serv *NewsServ) Get(news *NewsData) error {
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

//GetByArticleLink Params: news *model.NewsData; Returns: error; Description: Gets News Data piece by its article link
func (serv *NewsServ) GetByArticleLink(news *NewsData) error {
	vErr := validateID(news.ID)
	if vErr != nil {
		return vErr
	}
	alErr := validateString(news.ArticleLink, "Article Link")
	if alErr != nil {
		return alErr
	}
	err := serv.Repo.GetByArticleLink(news)
	if err != nil {
		return err
	}
	return nil
}
func (serv *NewsServ) GetNewsByWebNameAndID(news *NewsData) error {
	vErr := validateID(news.ID)
	if vErr != nil {
		return vErr
	}
	wnErr := validateString(news.WebsiteName, "Website Name")
	if wnErr != nil {
		return wnErr
	}
	err := serv.Repo.GetNewsByWebNameAndID(news)
	if err != nil {
		return err
	}
	return nil
}
func (serv *NewsServ) GetMultipleNews(start int, count int, news []*NewsData) error {
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
func (serv *NewsServ) GetMultipleNewsByWebName(start, count int, news []*NewsData, webName string) error {
	if start < 0 {
		err := errors.New("Start must be greater than zero.")
		return err
	}
	if count < 0 {
		err := errors.New("Start must be greater than zero.")
		return err
	}
	if len(news) < 0 {
		err := errors.New("Start must be greater than zero.")
		return err
	}

	err := serv.Repo.GetMultipleNewsByWebName(start, count, news, webName)

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
	newsDataSlice := []*NewsData{}
	for i := uint(1); i < processAmount; i++ {
		newsData := NewsDataInit(articleLink[i], articleText[i], imageURL[i], paragraphs[i-1], websiteName, i)
		newsDataSlice = append(newsDataSlice, newsData)
	}
	for i := 0; i < len(newsDataSlice); i++ {
		newsD := newsDataSlice[i]
		err := s.Repo.Create(newsD)
		if err != nil {
			s.Repo.RollBack()
			return err
		}
	}
	return nil

}
