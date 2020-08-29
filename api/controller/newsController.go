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

type NewsController struct {
	newsServ          *news.NewsServ
	scraperServiceMap map[string]domain.RawNewsService
}

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
		response := u.Message(false, "Invalid or empty form value: count")
		u.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	start, startErr := strconv.Atoi(r.FormValue("start"))
	if startErr != nil {
		response := u.Message(false, "Invalid or empty form value: start")
		u.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	news := []*domain.NewsData{}
	getErr := cntrl.newsServ.GetMultipleNewsByWebName(start, count, news, webName)
	if getErr != nil {
		response := u.Message(false, "Error Getting News Articles: please try again.")
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
		response := u.Message(false, "Could not find Article for given Webname and ID")
		u.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	newsArticle := domain.NewsDataNoInit()
	newsArticle.WebsiteName = webName
	newsArticle.ID = *uuidParse
	if getErr := cntrl.newsServ.GetNewsByWebNameAndID(newsArticle); getErr != nil {
		response := u.Message(false, "Could not find Article for given Webname and ID")
		u.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	u.RespondWithJSON(w, http.StatusOK, newsArticle)
}

func (cntrl *NewsController) PullNewsData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	webName := checkWebName(vars, w)
	scrapeServ := cntrl.FindScraperService(webName)
	rawData := scrapeServ.RawNewsScrape()
	//if rawData.
}

func (cntrl *NewsController) FindScraperService(webname string) domain.RawNewsService {
	scrapeServ := cntrl.scraperServiceMap[webname]
	if scrapeServ == nil {
		return nil
	}
	return scrapeServ
}

func checkWebName(vars map[string]string, w http.ResponseWriter) string {
	webName := vars["newsname"]

	if webName != "globalnews" {
		response := u.Message(false, "Invalid News Website Name")
		u.RespondWithError(w, http.StatusBadRequest, response)
	}

	return webName
}
