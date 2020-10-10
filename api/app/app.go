// app.go

package app

import (
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gorilla/mux"
	accServ "github.com/jmattson4/go-sample-api/account/service"
	c "github.com/jmattson4/go-sample-api/api/controller"
	mw "github.com/jmattson4/go-sample-api/api/middleware"
	newsServ "github.com/jmattson4/go-sample-api/news/service"
)

//App models the application.
type App struct {
	Router   *mux.Router
	AccServ  *accServ.AccountService
	NewsServ *newsServ.NewsServ
	Enforcer *casbin.Enforcer
}

//ConstructApp ...
func ConstructApp(enf *casbin.Enforcer, accserv *accServ.AccountService, newserv *newsServ.NewsServ) *App {
	app := &App{}
	app.Router = mux.NewRouter()
	app.AccServ = accserv
	app.NewsServ = newserv
	app.Enforcer = enf
	return app
}

//Initialize To be used before application is run in order to connect to the database and create routes.
func (a *App) Initialize() {
	a.initializeMiddleware()
	auth, n := a.createControllers()
	a.initializeRoutes(auth, n)
}

//InitializeTesting ...
//initialize To be used before application is run in order to connect to the database and create routes.
func (a *App) InitializeTesting() {
	auth, n := a.createControllers()
	a.initializeRoutes(auth, n)
}

func (a *App) createControllers() (*c.AccountController, *c.NewsController) {
	authController := c.ConstructAccountController(a.AccServ)
	newsController := c.ConstructNewsController(a.NewsServ)

	return authController, newsController
}

//InitializeRoutes to be used to create all the routes on the API
func (a *App) initializeRoutes(acc *c.AccountController, n *c.NewsController) {

	a.Router.HandleFunc("/api/news/{newsname}", n.GetNewsByWebName).Methods("GET")
	a.Router.HandleFunc("/api/news/{newsname}/{id}", n.GetNewsArticleByID).Methods("GET")
	a.Router.HandleFunc("/api/news/{newsname}/scrape", n.PullNewsData).Methods("POST")

	a.Router.HandleFunc("/api/user/new", acc.CreateAccount).Methods("POST")
	a.Router.HandleFunc("/api/user/login", acc.Authenticate).Methods("POST")
	a.Router.HandleFunc("/api/user/logout", acc.Logout).Methods("POST")
	a.Router.HandleFunc("/api/user/refresh", acc.Refresh).Methods("POST")
	a.Router.HandleFunc("/api/users", acc.GetAllAccounts).Methods("GET")
}

func (a *App) initializeMiddleware() {
	a.Router.Use(mw.JwtAuthentication(a.AccServ))
	a.Router.Use(mw.Authorize(a.Enforcer, a.AccServ))
}

//Run To be used to start up the server. Use after initilization.
func (a *App) Run(addr string) {
	log.Print("application starting on port 8010")
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
