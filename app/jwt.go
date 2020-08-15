package app

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jmattson4/go-sample-api/cache"
	u "github.com/jmattson4/go-sample-api/util"
)

type AccessDetails struct {
	AccessUuid string
	UserID     uint
}

//JwtAuthentication ... Handler to ensure that every
var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkEndpointList(next, w, r)
		response := make(map[string]interface{})
		tokenAuth, err := extractTokenMetaData(r)

		if err != nil {
			response = u.Message(false, "Token sent is unauthorized please login to get a new token.")
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
		fmt.Sprintf("User %", userID) //Useful for monitoring
		ctx := context.WithValue(r.Context(), "user", userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

func checkEndpointList(next http.Handler, w http.ResponseWriter, r *http.Request) {
	notAuth := []string{"/api/user/new", "/api/user/login"} //List of endpoints that doesn't require auth
	requestPath := r.URL.Path                               //current request path

	//check if request does not need authentication, serve the request if it doesn't need it
	for _, value := range notAuth {

		if value == requestPath {
			next.ServeHTTP(w, r)
			return
		}
	}
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
		return []byte(u.GetEnv().AccessSecret), nil
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
func extractTokenMetaData(r *http.Request) (*AccessDetails, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserID:     uint(userId),
		}, nil
	}
	return nil, err
}

func fetchAuth(auuthD *AccessDetails) (uint, error) {
	userid, err := cache.Client.Get(auuthD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return uint(userID), nil
}
