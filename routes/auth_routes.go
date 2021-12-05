package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
)

func AuthRoutes(router *mux.Router, db *sql.DB, prefix string){
	subRouter := router.PathPrefix(prefix).Subrouter()

	// handler_auth.SignupAuth()
	// subRouter.HandleFunc("/signup", handler_auth.SignupAuth)
}