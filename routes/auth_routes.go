package routes

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	handler_auth "github.com/v-inmar/goweb/handlers/auth_handlers"
)

func AuthRoutes(router *mux.Router, db *sql.DB, prefix string){
	subRouter := router.PathPrefix(prefix).Subrouter()

	subRouter.HandleFunc("/signup", handler_auth.SignupAuth(db)).Methods(http.MethodPost)
	subRouter.HandleFunc("/login", handler_auth.LoginAuth(db)).Methods(http.MethodPost)
}