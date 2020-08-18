package news

/*
	This File is the interface declartions for the news section of the application.
	It contains definitions for base Reading and Writing to external Datasources
	It also contains some other classes used for The News repository and Service that
	Connects to PostGreSQL. Can be expanded to fit other External Datasources.

*/
type Reader interface {
	Get(news *NewsData) error
	GetMultipleNews(start int, count int, news []*NewsData) error
}
type Writer interface {
	Create(news *NewsData) error
	Update(news *NewsData) error
	Delete(news *NewsData) error
}

type PqlBase interface {
	GetMultipleNewsByWebName(start, count int, news []*NewsData, webName string) error
	HardDelete(news *NewsData) error
	GetByArticleLink(news *NewsData) error
	GetNewsByWebNameAndID(news *NewsData) error
}

type BaseRepository interface {
	Reader
	Writer
}

type NewsRepository interface {
	BaseRepository
	PqlBase
	RollBack() error
}

type BaseService interface {
	Reader
	Writer
}

type NewsService interface {
	BaseService
	PqlBase
	ProcessNewsData(
		processAmount uint,
		websiteName string,
		articleLink []string,
		articleText []string,
		imageURL []string,
		paragraphs []string) error
}
