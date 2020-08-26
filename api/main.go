// main.go

package main

import (
	"log"

	accCache "github.com/jmattson4/go-sample-api/account/repository/cache"
	accPql "github.com/jmattson4/go-sample-api/account/repository/pql"
	accServ "github.com/jmattson4/go-sample-api/account/service"
	a "github.com/jmattson4/go-sample-api/api/app"
	sec "github.com/jmattson4/go-sample-api/api/security"
	c "github.com/jmattson4/go-sample-api/cache"
	db "github.com/jmattson4/go-sample-api/database"
	newsPql "github.com/jmattson4/go-sample-api/news/repository/pql"
	newsC "github.com/jmattson4/go-sample-api/news/repository/redis"
	newsServ "github.com/jmattson4/go-sample-api/news/service"
	"github.com/jmattson4/go-sample-api/util"
)

func main() {
	env := util.GetEnv()

	accountDb := db.InitAccountDB(env)
	accountCache := c.InitRedisCache(env, 0)

	newsDb := db.InitNewsDB(env)
	newsCache := c.InitRedisCache(env, 5)

	accPqlRepo := accPql.ConstructAccountsRepo(accountDb)
	accCacheRepo := accCache.ConstructAccountCacheRepo(accountCache)

	newsPqlRepo := newsPql.ConstructNewsRepo(newsDb)
	newsCacheRepo := newsC.ConstructCacheRepo(cache)

	accountService := accServ.ConstructAccountService(accPqlRepo, accCacheRepo)
	newsService := newsServ.ConstructNewsService(newsPqlRepo, newsCacheRepo)

	enforcer, err := sec.InitAuthorizationEnforcer()
	if err != nil {
		log.Fatal("Error Getting Auth Enforcer: %v", err)
	} else {
		a := a.ConstructApp(enforcer, accountService, newsService)

		a.Initialize()

		a.Run(":8010")
	}
}
