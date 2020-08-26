package cache

import (
	"fmt"

	"github.com/go-redis/redis/v7"

	"github.com/jmattson4/go-sample-api/util"
)

//InitRedisCache Creates a redis cache to be used throughout the Repository.
func InitRedisCache(env *util.Environmentals, dbNumber int) *redis.Client {
	hn := env.RedisHostname
	port := env.RedisPort
	dsn := fmt.Sprintf("%v:%v", hn, port)
	pw := env.RedisPassword

	if len(dsn) == 0 {
		dsn = "cache:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: pw,
		DB:       dbNumber,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connection to redis was a success!")
	}
	return client
}
