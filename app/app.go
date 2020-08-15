// app.go

package app

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gorilla/mux"
	c "github.com/jmattson4/go-sample-api/controller"
	mw "github.com/jmattson4/go-sample-api/middleware"
)

//App models the application.
type App struct {
	Router   *mux.Router
	DB       *sql.DB
	Enforcer *casbin.Enforcer
}

//Initialize To be used before application is run in order to connect to the database and create routes.
func (a *App) Initialize(enf *casbin.Enforcer) {

	a.Router = mux.NewRouter()
	a.Enforcer = enf

	a.initializeMiddleware()
	a.initializeRoutes()
}

//Initialize To be used before application is run in order to connect to the database and create routes.
func (a *App) InitializeTesting() {

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

//InitializeRoutes to be used to create all the routes on the API
func (a *App) initializeRoutes() {

	a.Router.HandleFunc("/api/news/{newsname}", c.GetNewsByWebName).Methods("GET")
	a.Router.HandleFunc("/api/news/{newsname}/{id:[0-9]+}", c.GetNewsArticleByID).Methods("GET")

	a.Router.HandleFunc("/api/products", c.GetProducts).Methods("GET")
	a.Router.HandleFunc("/api/product", c.CreateProduct).Methods("POST")
	a.Router.HandleFunc("/api/product/{id:[0-9]+}", c.GetProduct).Methods("GET")
	a.Router.HandleFunc("/api/product/{id:[0-9]+}", c.UpdateProduct).Methods("PUT")
	a.Router.HandleFunc("/api/product/{id:[0-9]+}", c.DeleteProduct).Methods("DELETE")
	a.Router.HandleFunc("/api/products/deleted", c.ShowDeletedProducts).Methods("GET")

	a.Router.HandleFunc("/api/user/new", c.CreateAccount).Methods("POST")
	a.Router.HandleFunc("/api/user/login", c.Authenticate).Methods("POST")
	a.Router.HandleFunc("/api/user/logout", c.Logout).Methods("POST")
	a.Router.HandleFunc("/api/user/refresh", c.Refresh).Methods("POST")
}

func (a *App) initializeMiddleware() {
	a.Router.Use(mw.JwtAuthentication)
	a.Router.Use(mw.Authorize(a.Enforcer))
}

//Run To be used to start up the server. Use after initilization.
func (a *App) Run(addr string) {
	log.Print("application starting on port 8010")
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
