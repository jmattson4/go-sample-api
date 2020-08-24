package accounts

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/jmattson4/go-sample-api/domain"
	u "github.com/jmattson4/go-sample-api/util"
	"github.com/twinj/uuid"
)

type CacheRepository struct {
	redis *redis.Client
}

//CreateToken : Function used to generate a access token that expires in 15 minutes and a Refresh Token that expries in 7 days
func (cr *CacheRepository) CreateToken(accountID uint) (*domain.TokenDetails, error) {
	accessSecret := u.GetEnv().AccessSecret
	refreshSecret := u.GetEnv().RefreshSecret
	//Create token details setting
	td := &domain.TokenDetails{}
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

func (cr *CacheRepository) CreateAuth(userid uint, td *domain.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := cr.redis.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := cr.redis.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func (cr *CacheRepository) DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := cr.redis.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
