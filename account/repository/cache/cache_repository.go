package accounts

import (
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/jmattson4/go-sample-api/domain"
)

//AccountCacheRepository ...
type AccountCacheRepository struct {
	Redis *redis.Client
}

//ConstructAccountCacheRepo ...
func ConstructAccountCacheRepo(redis *redis.Client) *AccountCacheRepository {
	return &AccountCacheRepository{
		Redis: redis,
	}
}

//GetAuth This returns the userUUID associated with the given accessUUID
func (cr *AccountCacheRepository) GetAuth(accessUUID string) (string, error) {
	userUUID, err := cr.Redis.Get(accessUUID).Result()
	if err != nil {
		return "", err
	}
	return userUUID, nil
}

//CreateAuth ...
func (cr *AccountCacheRepository) CreateAuth(accountUUID string, td *domain.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := cr.Redis.Set(td.AccessUUID, accountUUID, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := cr.Redis.Set(td.RefreshUUID, accountUUID, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

//DeleteAuth ...
func (cr *AccountCacheRepository) DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := cr.Redis.Del(givenUuid).Result()
	if err != nil {
		return 0, domain.ACCOUNT_CACHE_AUTH_DELETION
	}
	return deleted, nil
}
