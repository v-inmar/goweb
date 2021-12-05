package handler_auth

import (
	"database/sql"
	"log"
	"net/http"
)

func LoginAuth(db *sql.DB) http.HandlerFunc{
	return func(rw http.ResponseWriter, r *http.Request) {
		// needs implementation
		log.Fatal("Implementation needed")
	}
}