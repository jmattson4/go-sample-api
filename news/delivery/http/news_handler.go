package http

import (
	"net/http"
	"strconv"

	"github.com/twinj/uuid"

	"github.com/gorilla/mux"
	"github.com/jmattson4/go-sample-api/domain"
	"github.com/jmattson4/go-sample-api/util"

	u "github.com/jmattson4/go-sample-api/util"
)

type NewsHandler struct {
	NewsServ domain.NewsService
}

func ConstructNewsHandler(router *mux.Router, service domain.NewsService) *NewsHandler {
	newsHandler := &NewsHandler{
		NewsServ: service,
	}
	router.HandleFunc(domain.NEWS_ROUTE_GET_BY_WEBNAME, newsHandler.GetNewsByWebName).Methods("GET")
	router.HandleFunc(domain.NEWS_ROUTE_GET_BY_WEBNAME, GetNewsByWebName).Methods("GET")
	return newsHandler
}

//GetNewsByWebName ...
//Description: This function is a Route that takes the news Website Name, a start and a count value.
// Then it Checks against scrapped website names finally displaying all articles within the start and count interval
func (hand *NewsHandler) GetNewsByWebName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	webName := checkWebName(vars, w)
	count, countErr := strconv.Atoi(r.FormValue("count"))
	if countErr != nil {
		response := util.Message(false, "Invalid or empty form value: count")
		u.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	start, startErr := strconv.Atoi(r.FormValue("start"))
	if startErr != nil {
		response := util.Message(false, "Invalid or empty form value: start")
		u.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	news := &[]m.NewsData{}
	getErr := m.GetMultipleNewsByWebName(start, count, news, webName)
	if getErr != nil {
		response := util.Message(false, "Error Getting News Articles: please try again.")
		u.RespondWithError(w, http.StatusInternalServerError, response)
		return
	}

	u.RespondWithJSON(w, http.StatusOK, news)
}

//GetNewsArticleByID ...
//Description: This function finds a specific article by its webname and ID
func (hand *NewsHandler) GetNewsArticleByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	webName := checkWebName(vars, w)
	id := vars["id"]
	uuidParse, uuidErr := uuid.Parse(id)
	if uuidErr != nil {
		response := util.Message(false, "Could mot find Article for given Webname and ID")
		util.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	newsArticle := domain.NewsDataNoInit()
	newsArticle.WebsiteName = webName
	newsArticle.ID = *uuidParse
	if getErr := hand.NewsServ.GetNewsByWebNameAndID(newsArticle); getErr != nil {
		response := util.Message(false, "Could not find Article for given Webname and ID")
		util.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	u.RespondWithJSON(w, http.StatusOK, newsArticle)
}

func checkWebName(vars map[string]string, w http.ResponseWriter) string {
	webName := vars["newsname"]

	if webName != "globalnews" {
		response := util.Message(false, "Invalid News Website Name")
		u.RespondWithError(w, http.StatusBadRequest, response)
	}

	return webName
}
