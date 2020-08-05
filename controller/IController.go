package controller

import "database/sql"

//IController interface for a Http handler/controller
type IController interface {
	InitController(*sql.DB)
}
