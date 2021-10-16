package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	handlers "github.com/v-inmar/goweb/handlers/errors"
	m "github.com/v-inmar/goweb/models"
)

func main() {
	// Initialize a new mux router
	router := mux.NewRouter()

	// Assign a handle function for the route '/test'
	// to test the connection and router is working
	// Anonymous function is used to handle this route
	// Declared as GET
	// Will default to 200 OK status code
	router.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
		// Set the header
		rw.Header().Set("Content-Type", "application/json")

		// set the payload and response data
		rm := m.ResponseModel{
			Status: "200 OK",
			Payload: m.PayloadModel{
				Payload: "API is working",
			},
		}

		// JSON encode the connection value and pass to ResponseWriter
		json.NewEncoder(rw).Encode(rm)

	}).Methods("GET")

	// Assign the NotFoundHandler (custom) to mux's NotFoundHandler
	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	// Run server with the router and log error if failed
	log.Fatal(http.ListenAndServe(":5000", router))

}
