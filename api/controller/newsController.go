package controller

import (
	"net/http"
	"strconv"

	m "github.com/jmattson4/go-sample-api/model"
	"github.com/jmattson4/go-sample-api/util"

	"github.com/gorilla/mux"
	u "github.com/jmattson4/go-sample-api/util"
)

//GetNewsByWebName ...
//Description: This function is a Route that takes the news Website Name, a start and a count value.
// Then it Checks against scrapped website names finally displaying all articles within the start and count interval
func GetNewsByWebName(w http.ResponseWriter, r *http.Request) {
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
func GetNewsArticleByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	webName := checkWebName(vars, w)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response := util.Message(false, "Invalid News Article ID")
		util.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	newsArticle := &m.NewsData{}
	newsArticle.WebsiteName = webName
	newsArticle.ID = uint(id)
	if getErr := newsArticle.GetNewsByWebNameAndID(); getErr != nil {
		response := util.Message(false, "Could mot find Article for given Webname and ID")
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
