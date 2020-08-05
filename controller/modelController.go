package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	m "github.com/jmattson4/go-sample-api/model"
	"github.com/jmattson4/go-sample-api/util"
)

//ModelController ...
type ModelController struct {
	DB *sql.DB
}

//InitController ...
func (mc *ModelController) InitController(db *sql.DB) {
	mc.DB = db
}

// GetProduct ...
func (mc *ModelController) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response := util.Message(false, "Invalid product ID")
		util.RespondWithError(w, http.StatusBadRequest, response)
		return
	}

	p := m.Product{ID: id}
	if err := p.GetProduct(mc.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			response := util.Message(false, "Invalid product ID")
			util.RespondWithError(w, http.StatusNotFound, response)
		default:
			response := util.Message(false, err.Error())
			util.RespondWithError(w, http.StatusInternalServerError, response)
		}
		return
	}

	util.RespondWithJSON(w, http.StatusOK, p)
}

//GetProducts ...
func (mc *ModelController) GetProducts(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := m.GetProducts(mc.DB, start, count)
	if err != nil {
		response := util.Message(false, err.Error())
		util.RespondWithError(w, http.StatusInternalServerError, response)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, products)
}

//CreateProduct ...
func (mc *ModelController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p m.Product
	//creates decodeer based on the request body
	decoder := json.NewDecoder(r.Body)
	//decodes the json body from the request into the newly created product struct.
	//	if there is errors with the decoding the route responds with a HTTP 400 response
	if err := decoder.Decode(&p); err != nil {
		response := util.Message(false, "Invalid request payload")
		util.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()
	//attempts to create the product which was sent via json into the
	// database. If there is an error it throws an HTTP 500 error
	if err := p.CreateProduct(mc.DB); err != nil {
		response := util.Message(false, "Invalid request payload")
		util.RespondWithError(w, http.StatusInternalServerError, response)
		return
	}
	//if everything is okay then api responds with 201 status created.
	//	and sends the newly created product back to the user.
	util.RespondWithJSON(w, http.StatusCreated, p)
}

//UpdateProducts ...
func (mc *ModelController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response := util.Message(false, "Invalid request payload")
		util.RespondWithError(w, http.StatusBadRequest, response)
		return
	}

	var p m.Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		response := util.Message(false, "Invalid request payload")
		util.RespondWithError(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()
	p.ID = id

	if err := p.UpdateProduct(mc.DB); err != nil {
		response := util.Message(false, err.Error())
		util.RespondWithError(w, http.StatusInternalServerError, response)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, p)
}

//DeleteProducts ...
func (mc *ModelController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response := util.Message(false, "Invalid Product ID")
		util.RespondWithError(w, http.StatusBadRequest, response)
		return
	}

	p := m.Product{ID: id}
	if err := p.DeleteProduct(mc.DB); err != nil {
		response := util.Message(false, err.Error())
		util.RespondWithError(w, http.StatusInternalServerError, response)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
