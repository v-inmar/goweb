package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	errorhandlers "github.com/v-inmar/goweb/handlers/errors"

	todo_routes "github.com/v-inmar/goweb/routes"
)

func main() {

	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Create instance of the App struct
	app := App{}

	// golang run main.go db make
	// a way of auto creating the database using the given sql file
	if len(os.Args) > 1 {
		if os.Args[1] == "db" {
			if os.Args[2] == "make" {
				app.MakeDB(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBHOST"), os.Getenv("DBPORT"))
				os.Exit(1)
			}
		}
	}

	// Initliaze the app with connection to the db
	app.Initialize(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBNAME"))

	// Initialize a new mux router
	// router := mux.NewRouter()
	router := app.Router // Using the mux router within the app.go

	// Assign a handle function for the route '/test'
	// to test the connection and router is working
	// Anonymous function is used to handle this route
	// Declared as GET
	// Will default to 200 OK status code
	router.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
		// Set the header
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
	}).Methods("GET")

	// Create todos sub routers
	todo_routes.TodoRoutes(app.Router, app.DB, "/todos")

	// Assign the NotFoundHandler (custom) to mux's NotFoundHandler
	router.NotFoundHandler = http.HandlerFunc(errorhandlers.NotFoundHandler)

	// Run server with the router and log error if failed
	// log.Fatal(http.ListenAndServe(":5000", router))
	app.Run(":5000") // app.go AppRun that runs the server

}
