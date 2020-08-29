package news

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/jmattson4/go-sample-api/domain"
)

type NewsCacheRepo struct {
	Redis *redis.Client
}

func ConstructCacheRepo(redis *redis.Client) *NewsCacheRepo {
	return &NewsCacheRepo{
		Redis: redis,
	}
}
func generateExpirationTime() time.Duration {
	expiresAt := time.Now().Add(time.Hour * 24).Unix()
	et := time.Unix(expiresAt, 0)
	now := time.Now()
	return et.Sub(now)
}

//Create Params: news *model.NewsData; Returns: error; Description: Creates Scraped News Data piece
func (repo *NewsCacheRepo) Create(news *domain.NewsData) error {
	encode, encodeErr := json.Marshal(news)
	if encodeErr != nil {
		return encodeErr
	}
	expiration := generateExpirationTime()
	err := repo.Redis.Set(fmt.Sprintf("id:%v", news.ID), encode, expiration)
	if err.Err() != nil {
		return err.Err()
	}
	return nil
}

//Update Params: news *model.NewsData; Returns: error; Description: Updates News Data piece
func (repo *NewsCacheRepo) Update(news *domain.NewsData) error {
	encode, encodeErr := json.Marshal(news)
	if encodeErr != nil {
		return encodeErr
	}
	expiration := generateExpirationTime()
	err := repo.Redis.Set(fmt.Sprintf("id:%v", news.ID), encode, expiration)
	if err.Err() != nil {
		return err.Err()
	}
	return nil
}
func (repo *NewsCacheRepo) Delete(news *domain.NewsData) error {
	_, err := repo.Redis.Del(fmt.Sprintf("id:%v", news.ID)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (repo *NewsCacheRepo) HardDelete(news *domain.NewsData) error {
	_, err := repo.Redis.Del(fmt.Sprintf("id:%v", news.ID)).Result()
	if err != nil {
		return err
	}
	return nil
}

//Get Params: news *model.NewsData; Returns: error; Description: Gets News Data piece
func (repo *NewsCacheRepo) Get(news *domain.NewsData) error {
	value, err := repo.Redis.Get(fmt.Sprintf("id:%v", news.ID)).Result()
	if err != nil {
		return err
	}
	unMarhsalErr := json.Unmarshal([]byte(value), news)
	if unMarhsalErr != nil {
		return unMarhsalErr
	}
	return err
}
