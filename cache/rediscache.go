package cache

import (
	"fmt"

	"github.com/go-redis/redis/v7"

	"github.com/jmattson4/go-sample-api/util"
)

var Client *redis.Client

func init() {
	hn := util.GetEnv().RedisHostname
	port := util.GetEnv().RedisPort
	dsn := fmt.Sprintf("%v:%v", hn, port)
	pw := util.GetEnv().RedisPassword

	if len(dsn) == 0 {
		dsn = "cache:6379"
	}
	Client = redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: pw,
	})
	_, err := Client.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connection to redis was a success!")
	}
}
