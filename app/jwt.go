package app

import (
	"net/http"
	u "github.com/jmattson4/go-sample-api/util"
	"strings"
	jwt "github.com/dgrijalva/jwt-go"
	"os"
	"context"
	"fmt"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {

	notAuth := []string{"/api/user/new", "/api/user/login"} //List of endpoints that doesn't require auth
	requestPath := r.URL.Path //current request path

	//check if request does not need authentication, serve the request if it doesn't need it
	for _, value := range notAuth {

		if value == requestPath {
			next.ServeHTTP(w, r)
			return
		}
	}

	response := make(map[string] interface{})
	tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

	if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
		response = u.Message(false, "Missing auth token")
		w.Header().Add("Content-Type", "application/json")	
		u.RespondWithError(w, http.StatusForbidden, response)
		return
	}

	splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	if len(splitted) != 2 {
		response = u.Message(false, "Invalid/Malformed auth token")
		w.Header().Add("Content-Type", "application/json")
		u.RespondWithError(w, http.StatusForbidden, response)
		return
	}

	tokenPart := splitted[1] //Grab the token part, what we are truly interested in
	tk := &models.Token{}

	token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})

	if err != nil { //Malformed token, returns with http code 403 as usual
		response = u.Message(false, "Malformed authentication token")
		w.Header().Add("Content-Type", "application/json")
		u.RespondWithError(w, http.StatusForbidden, response)
		return
	}

	if !token.Valid { //Token is invalid, maybe not signed on this server
		response = u.Message(false, "Token is not valid.")
		w.Header().Add("Content-Type", "application/json")
		u.RespondWithError(w, http.StatusForbidden, response)
		return
	}

	//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
	fmt.Sprintf("User %", tk.Username) //Useful for monitoring
	ctx := context.WithValue(r.Context(), "user", tk.UserId)
	r = r.WithContext(ctx)
	next.ServeHTTP(w, r) //proceed in the middleware chain!
});