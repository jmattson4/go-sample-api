package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	u "github.com/jmattson4/go-sample-api/api/utils"
	"github.com/jmattson4/go-sample-api/domain"
	news "github.com/jmattson4/go-sample-api/news/service"
	"github.com/twinj/uuid"
)

//NewsController ...
//	Description: Data Structure which stores the services used
//	by the News Controller
type NewsController struct {
	newsServ          *news.NewsServ
	scraperServiceMap map[string]domain.RawNewsService
}

//ConstructNewsController ...
//	Description: Constructor for the News Controller
func ConstructNewsController(newserv *news.NewsServ) *NewsController {
	return &NewsController{
		newsServ: newserv,
	}
}

//GetNewsByWebName ...
//Description: This function is a Route that takes the news Website Name, a start and a count value.
// Then it Checks against scrapped website names finally displaying all articles within the start and count interval
func (cntrl *NewsController) GetNewsByWebName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	webName := checkWebName(vars, w)
	count, countErr := strconv.Atoi(r.FormValue("count"))
	if countErr != nil {
		response := u.Message(false, domain.NEWS_CONTROLLER_GETNEWSBYWEBNAME_INVALID_COUNT_FIELD)
		u.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	start, startErr := strconv.Atoi(r.FormValue("start"))
	if startErr != nil {
		response := u.Message(false, domain.NEWS_CONTROLLER_GETNEWSBYWEBNAME_INVALID_START_FIELD)
		u.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	news := []*domain.NewsData{}
	getErr := cntrl.newsServ.GetMultipleNewsByWebName(start, count, news, webName)
	if getErr != nil {
		response := u.Message(false, domain.NEWS_CONTROLLER_GETNEWSBYWEBNAME_INVALID_GET_ERROR)
		u.RespondWithError(w, http.StatusInternalServerError, response)
		return
	}

	u.RespondWithJSON(w, http.StatusOK, news)
}

//GetNewsArticleByID ...
//Description: This function finds a specific article by its webname and ID
func (cntrl *NewsController) GetNewsArticleByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	webName := checkWebName(vars, w)
	id := vars["id"]
	uuidParse, uuidErr := uuid.Parse(id)
	if uuidErr != nil {
		response := u.Message(false, domain.NEWS_CONTROLLER_GETNEWSARTICLEBYID_UUID_PARSE)
		u.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	newsArticle := domain.NewsDataNoInit()
	newsArticle.WebsiteName = webName
	newsArticle.ID = *uuidParse
	if getErr := cntrl.newsServ.GetNewsByWebNameAndID(newsArticle); getErr != nil {
		response := u.Message(false, domain.NEWS_CONTROLLER_GETNEWSARTICLEBYID_ARTICLE_NOT_FOUND)
		u.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	u.RespondWithJSON(w, http.StatusOK, newsArticle)
}

//PullNewsData ...
//	Description: This function will initate a scrape via scrape service and then process and save the data
//	into the database via the newsService
func (cntrl *NewsController) PullNewsData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	webName := checkWebName(vars, w)
	scrapeServ := cntrl.FindScraperService(webName)
	rawNewsData, scrapeErr := scrapeServ.RawNewsScrape()
	//if rawData.
	if scrapeErr != nil {
		response := u.Message(false, domain.NEWS_CONTROLLER_PULLNEWSDATA_FAILED_SCRAPE)
		u.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	processErr := cntrl.newsServ.ProcessNewsData(20, webName, rawNewsData)
	if processErr != nil {
		response := u.Message(false, domain.NEWS_CONTROLLER_PULLNEWSDATA_SUCCESS)
		u.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	u.RespondWithJSON(w, http.StatusOK, domain.NEWS_CONTROLLER_PULLNEWSDATA_SUCCESS)
}

//FindScraperService ...
//	Gets the proper ScraperService for the specific website.
func (cntrl *NewsController) FindScraperService(webname string) domain.RawNewsService {
	scrapeServ := cntrl.scraperServiceMap[webname]
	if scrapeServ == nil {
		return nil
	}
	return scrapeServ
}

// checkWebName
//	 This function checks the newsname requested in a GET and checks it
//	 against the list of possible options. It returns an error message if they chose
//	 an incorrect option.
func checkWebName(vars map[string]string, w http.ResponseWriter) string {
	webName := vars["newsname"]

	if webName != "globalnews" {
		response := u.Message(false, domain.NEWS_CONTROLLER_CHECKWEBNAME_INVALID_WEBNAME)
		u.RespondWithError(w, http.StatusBadRequest, response)
	}

	return webName
}
