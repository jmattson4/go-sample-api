// main.go

package main

import (
	"fmt"
	"log"

	accCache "github.com/jmattson4/go-sample-api/account/repository/cache"
	accPql "github.com/jmattson4/go-sample-api/account/repository/pql"
	accServ "github.com/jmattson4/go-sample-api/account/service"
	a "github.com/jmattson4/go-sample-api/api/app"
	c "github.com/jmattson4/go-sample-api/cache"
	db "github.com/jmattson4/go-sample-api/database"
	newsPql "github.com/jmattson4/go-sample-api/news/repository/pql"
	newsC "github.com/jmattson4/go-sample-api/news/repository/redis"
	newsServ "github.com/jmattson4/go-sample-api/news/service"
	sec "github.com/jmattson4/go-sample-api/security"
	"github.com/jmattson4/go-sample-api/util"
)

func main() {
	util.ConstructEnv()
	env := util.GetEnv()

	accountDb := db.InitAccountDB(
		env.AccountDBService,
		env.AccountDBPort,
		env.AccountUser,
		env.AccountPassword,
		env.AccountDatabaseName)
	accountCache := c.InitRedisCache(env, 0)

	newsDb := db.InitNewsDB(
		env.DatabaseUser,
		env.DatabasePassword,
		env.DatabaseName,
		env.DatabaseDBService,
		env.DatabaseDBPort)
	newsCache := c.InitRedisCache(env, 5)

	accPqlRepo := accPql.ConstructAccountsRepo(accountDb)
	accCacheRepo := accCache.ConstructAccountCacheRepo(accountCache)

	newsPqlRepo := newsPql.ConstructNewsRepo(newsDb)
	newsCacheRepo := newsC.ConstructCacheRepo(newsCache)

	accountService := accServ.ConstructAccountService(accPqlRepo, accCacheRepo)
	newsService := newsServ.ConstructNewsService(newsPqlRepo, newsCacheRepo)

	enforcer, err := sec.InitAuthorizationEnforcer()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error Getting Auth Enforcer: %v", err))
	} else {
		a := a.ConstructApp(enforcer, accountService, newsService)

		a.Initialize()

		a.Run(":8010")
	}
}
