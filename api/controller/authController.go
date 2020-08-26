package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	acc "github.com/jmattson4/go-sample-api/account/service"
	u "github.com/jmattson4/go-sample-api/api/utils"
	"github.com/jmattson4/go-sample-api/domain"
	"github.com/jmattson4/go-sample-api/util"
)

type AuthController struct {
	accServ *acc.AccountService
}

func ConstructAuthController(acc *acc.AccountService) *AuthController {
	return &AuthController{
		accServ: acc,
	}
}

//CreateAccount ...
func (auth *AuthController) CreateAccount(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")

	err := auth.accServ.Create(email, password) //Create account
	if err != nil {
		u.RespondWithError(w, http.StatusForbidden, u.Message(false, fmt.Sprintf("Error: %v", err.Error())))
		return
	}
	u.RespondWithJSON(w, http.StatusOK, u.Message(false, domain.ACCOUNT_CREATION_SUCCESS))
}

//Authenticate ...
func (auth *AuthController) Authenticate(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")

	acc, err := auth.accServ.Login(email, password)
	if err != nil {
		u.RespondWithError(w, http.StatusForbidden, u.Message(false, fmt.Sprintf("Error: %v", err.Error())))
		return
	}
	u.RespondWithJSON(w, http.StatusOK, acc)
}

//Logout used to logout. Deleted the stored access token in the redis cache
func (auth *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	au, err := u.ExtractTokenMetaData(r)
	if err != nil {
		u.RespondWithError(w, http.StatusForbidden, u.Message(false, "Unauthorized"))
		return
	}
	delErr := auth.accServ.Logout(au.AccessUuid)
	if delErr != nil {
		u.RespondWithError(w, http.StatusForbidden, u.Message(false, "Cannot Deleted: Unauthorized"))
		return
	}
	u.RespondWithJSON(w, http.StatusOK, "Successfully logged out!")
}

//Refresh used to refresh the current refresh token gives back a new refresh and access
func (auth *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
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
		return []byte(util.GetEnv().RefreshSecret), nil
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
		userId, err := claims["account_id"].(string)
		if !err {
			u.RespondWithError(w, http.StatusUnprocessableEntity, u.Message(false, "Error occured"))
			return
		}
		//Delete the previous Refresh Token
		delErr := auth.accServ.DeleteAuth(refreshUuid)
		if delErr != nil { //if any goes wrong
			u.RespondWithError(w, http.StatusUnauthorized, u.Message(false, "Unauthorized"))
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := auth.accServ.CreateToken(userId)
		if createErr != nil {
			u.RespondWithError(w, http.StatusForbidden, u.Message(false, createErr.Error()))
			return
		}
		//save the tokens metadata to redis
		saveErr := auth.accServ.CreateAuth(userId, ts)
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
