package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	u "github.com/jmattson4/go-sample-api/api/utils"
	"github.com/jmattson4/go-sample-api/cache"
	"github.com/jmattson4/go-sample-api/domain"
	"github.com/jmattson4/go-sample-api/util"
)

type AccessDetails struct {
	AccessUuid string
	UserID     string
}

//JwtAuthentication ... Handler to ensure that every
var JwtAuthentication = func(next http.Handler) http.Handler {
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
		tokenAuth, err := ExtractTokenMetaData(r)

		if err != nil {
			response = u.Message(false, fmt.Sprintf("Token sent is unauthorized please login to get a new token: %v", err))
			u.RespondWithError(w, http.StatusUnauthorized, response)
			return
		}
		userID, fetchErr := fetchAuth(tokenAuth)

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

func extractToken(r *http.Request) string {
	tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
	splitted := strings.Split(tokenHeader, " ")  //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	if len(splitted) == 2 {
		return splitted[1]
	}
	return ""
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//make sure that token signing method conform to SigningMethodHMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(util.GetEnv().AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func tokenValid(r *http.Request) error {
	token, err := verifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

//ExtractTokenMetaData This function is used to pull the access token meta data from a request
//	in order to check the cache.
func ExtractTokenMetaData(r *http.Request) (*AccessDetails, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, domain.JWT_CANNOT_FIND_PROPERTY_ACCESS
		}
		userId, ok2 := claims["account_id"].(string)
		if !ok2 {
			return nil, domain.JWT_CANNOT_FIND_PROPERTY_ACCOUNT
		}

		return &AccessDetails{
			AccessUuid: accessUuid,
			UserID:     userId,
		}, nil
	}
	return nil, err
}

func fetchAuth(auuthD *AccessDetails) (string, error) {
	userid, err := cache.Client.Get(auuthD.AccessUuid).Result()
	if err != nil {
		return "", err
	}
	return userid, nil
}
