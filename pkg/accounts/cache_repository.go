package accounts

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmattson4/go-sample-api/cache"
	u "github.com/jmattson4/go-sample-api/util"
	"github.com/twinj/uuid"
)

//CreateToken : Function used to generate a access token that expires in 15 minutes and a Refresh Token that expries in 7 days
func CreateToken(accountID uint) (*TokenDetails, error) {
	accessSecret := u.GetEnv().AccessSecret
	refreshSecret := u.GetEnv().RefreshSecret
	//Create token details setting
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["account_id"] = accountID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), atClaims)
	td.AccessToken, err = at.SignedString([]byte(accessSecret))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["account_id"] = accountID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(refreshSecret))
	if err != nil {
		return nil, err
	}

	return td, nil

}

func CreateAuth(userid uint, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := cache.Client.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := cache.Client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := cache.Client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
