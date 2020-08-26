package app

import (
	"context"
	"fmt"
	"net/http"

	accServ "github.com/jmattson4/go-sample-api/account/service"
	u "github.com/jmattson4/go-sample-api/api/utils"
)

//JwtAuthentication ... Handler to ensure that every
func JwtAuthentication(serv *accServ.AccountService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			notAuth := []string{"/api/user/new", "/api/user/login", "/api/user/refresh"} //List of endpoints that doesn't require auth
			requestPath := r.URL.Path                                                    //current request path

			//check if request does not need authentication, serve the request if it doesn't need it
			for _, value := range notAuth {

				if value == requestPath {
					next.ServeHTTP(w, r)
					return
				}
			}

			response := make(map[string]interface{})
			tokenAuth, err := u.ExtractTokenMetaData(r)

			if err != nil {
				response = u.Message(false, fmt.Sprintf("Token sent is unauthorized please login to get a new token: %v", err))
				u.RespondWithError(w, http.StatusUnauthorized, response)
				return
			}
			userID, fetchErr := serv.GetAuth(tokenAuth.AccessUuid)

			if fetchErr != nil {
				response = u.Message(false, "Could not fetch User. Please login to get a new token.")
				u.RespondWithError(w, http.StatusUnauthorized, response)
				return
			}

			//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
			//fmt.Sprintf("User %v", userID) //Useful for monitoring
			ctx := context.WithValue(r.Context(), "user", userID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r) //proceed in the middleware chain!
		})
	}
}
