package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	errorhandlers "github.com/v-inmar/goweb/handlers/errors"
	handlers "github.com/v-inmar/goweb/handlers/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	theApp := App{}
	port, err := strconv.Atoi(os.Getenv("DBPORT"))
	if err != nil {
		log.Fatal(err)
	}
	theApp.AppInit(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBHOST"), port, os.Getenv("DBNAME"))

	// Initialize a new mux router
	// router := mux.NewRouter()

	// Assign a handle function for the route '/test'
	// to test the connection and router is working
	// Anonymous function is used to handle this route
	// Declared as GET
	// Will default to 200 OK status code
	theApp.Router.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
		// Set the header
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
	}).Methods("GET")

	// Get all todos
	theApp.Router.HandleFunc("/todos", handlers.GetAllTodos).Methods(http.MethodGet)

	// Create todo
	theApp.Router.HandleFunc("/todos", handlers.CreateTodo).Methods(http.MethodPost)

	// Get todo by id
	theApp.Router.HandleFunc("/todos/{id}", handlers.GetTodo).Methods(http.MethodGet)

	// Update todo by id
	theApp.Router.HandleFunc("/todos/{id}", handlers.UpdateTodo).Methods(http.MethodPut)

	// Delete a todo item by id
	theApp.Router.HandleFunc("/todos/{id}", handlers.DeleteTodo).Methods(http.MethodDelete)

	// Assign the NotFoundHandler (custom) to mux's NotFoundHandler
	theApp.Router.NotFoundHandler = http.HandlerFunc(errorhandlers.NotFoundHandler)

	// Run server with the router and log error if failed
	// log.Fatal(http.ListenAndServe(":5000", router))
	theApp.AppRun(":5000")

}
