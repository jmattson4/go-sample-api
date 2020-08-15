package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	mw "github.com/jmattson4/go-sample-api/middleware"
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

//Logout used to logout. Deleted the stored access token in the redis cache
var Logout = func(w http.ResponseWriter, r *http.Request) {
	au, err := mw.ExtractTokenMetaData(r)
	if err != nil {
		u.RespondWithError(w, http.StatusForbidden, u.Message(false, "Unauthorized"))
		return
	}
	deleted, delErr := model.DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 {
		u.RespondWithError(w, http.StatusForbidden, u.Message(false, "Cannot Deleted: Unauthorized"))
		return
	}
	u.RespondWithJSON(w, http.StatusOK, "Successfully logged out!")
}

//Refresh used to refresh the current refresh token gives back a new refresh and access
var Refresh = func(w http.ResponseWriter, r *http.Request) {
	mapToken := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&mapToken); err != nil {
		u.RespondWithError(w, http.StatusUnprocessableEntity, u.Message(false, "Invalid Map token"))
		return
	}
	refreshToken := mapToken["refresh_token"]
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(u.GetEnv().RefreshSecret), nil
	})
	if err != nil {
		u.RespondWithError(w, http.StatusUnauthorized, u.Message(false, "Refresh Token has expired"))
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		u.RespondWithError(w, http.StatusUnauthorized, u.Message(false, fmt.Sprintf("%v", err)))
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			u.RespondWithError(w, http.StatusUnprocessableEntity, u.Message(false, fmt.Sprintf("%v", err)))
			return
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["account_id"]), 10, 64)
		if err != nil {
			u.RespondWithError(w, http.StatusUnprocessableEntity, u.Message(false, "Error occured"))
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := model.DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			u.RespondWithError(w, http.StatusUnauthorized, u.Message(false, "Unauthorized"))
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := model.CreateToken(uint(userId))
		if createErr != nil {
			u.RespondWithError(w, http.StatusForbidden, u.Message(false, createErr.Error()))
			return
		}
		//save the tokens metadata to redis
		saveErr := model.CreateAuth(uint(userId), ts)
		if saveErr != nil {
			u.RespondWithError(w, http.StatusForbidden, u.Message(false, saveErr.Error()))
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		u.RespondWithJSON(w, http.StatusCreated, tokens)
	} else {
		u.RespondWithJSON(w, http.StatusUnauthorized, "refresh expired")
	}
}
