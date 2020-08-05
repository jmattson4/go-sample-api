// app.go

package app

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	mc "github.com/jmattson4/go-sample-api/controller"
	db "github.com/jmattson4/go-sample-api/database"
	_ "github.com/lib/pq"
)

//App models the application.
type App struct {
	Router          *mux.Router
	DB              *sql.DB
	ModelController *mc.ModelController
}

//Initialize To be used before application is run in order to connect to the database and create routes.
func (a *App) Initialize(user string, password string, dbname string, instanceConnectionName string) {

	a.Router = mux.NewRouter()
	var err error
	a.DB, err = db.InitDatabase(user, password, dbname, instanceConnectionName)
	if err != nil {
		log.Fatal("Database Error: Initialization failed.")
		defer a.DB.Close()
		return
	}
	a.ModelController = mc.InitModelController(a.DB)

	a.initializeRoutes()
}

//InitializeRoutes to be used to create all the routes on the API
func (a *App) initializeRoutes() {

	a.Router.HandleFunc("/products", a.ModelController.GetProduct).Methods("GET")
	a.Router.HandleFunc("/product", a.ModelController.CreateProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.ModelController.GetProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.ModelController.UpdateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.ModelController.DeleteProduct).Methods("DELETE")
}

//Run To be used to start up the server. Use after initilization.
func (a *App) Run(addr string) {
	log.Print("application starting on port 8010")
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}
