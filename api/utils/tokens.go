package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmattson4/go-sample-api/domain"
	"github.com/jmattson4/go-sample-api/util"
)

type AccessDetails struct {
	AccessUuid string
	UserID     string
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
