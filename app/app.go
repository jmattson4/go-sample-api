// app.go

package app

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	c "github.com/jmattson4/go-sample-api/controller"
	_ "github.com/lib/pq"
)

//App models the application.
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

//Initialize To be used before application is run in order to connect to the database and create routes.
func (a *App) Initialize() {

	a.Router = mux.NewRouter()
	var err error
	if err != nil {
		log.Fatal("Database Error: Initialization failed.")
		defer a.DB.Close()
		return
	}

	a.initializeRoutes()
}

//InitializeRoutes to be used to create all the routes on the API
func (a *App) initializeRoutes() {

	a.Router.Use(JwtAuthentication)
	a.Router.HandleFunc("/products", c.GetProducts).Methods("GET")
	a.Router.HandleFunc("/product", c.CreateProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", c.GetProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", c.UpdateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", c.DeleteProduct).Methods("DELETE")
	a.Router.HandleFunc("/products/deleted", c.ShowDeletedProducts).Methods("GET")

	a.Router.HandleFunc("/api/user/new", c.CreateAccount).Methods("POST")
	a.Router.HandleFunc("/api/user/login", c.Authenticate).Methods("POST")
}

//Run To be used to start up the server. Use after initilization.
func (a *App) Run(addr string) {
	log.Print("application starting on port 8010")
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}
