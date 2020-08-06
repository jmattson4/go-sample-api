package controller

import (
	"encoding/json"
	"net/http"

	"github.com/jmattson4/go-sample-api/model"
	u "github.com/jmattson4/go-sample-api/util"
)

//CreateAccount ...
var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &model.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.RespondWithError(w, http.StatusForbidden, u.Message(false, "Invalid request"))
		return
	}
	defer r.Body.Close()
	resp := account.Create() //Create account
	u.RespondWithJSON(w, http.StatusOK, resp)
}

//Authenticate ...
var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &model.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.RespondWithError(w, http.StatusForbidden, u.Message(false, "Invalid request"))
		return
	}

	resp := model.Login(account.Email, account.Password)
	u.RespondWithJSON(w, http.StatusOK, resp)
}
