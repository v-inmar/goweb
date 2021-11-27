package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	DB *sql.DB
	Router *mux.Router
}


/*
Establishes connection with the database
and creates new mux router
assigning variables to the App struct fields
*/
func (app *App) Initialize(dbuser string, dbpass string, dbhost string, dbport string, dbname string){
	cfg := mysql.Config{
		User: dbuser,
		Passwd: dbpass,
		Net: "tcp",
		Addr: dbhost+":"+dbport,
		DBName: dbname,
		ParseTime: true, // matches Date and Datetime to go time.Time
	}

	var err error
	app.DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := app.DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected to ", dbname, " located at ", dbhost, ":", dbport)
	app.Router = mux.NewRouter()
}

func (app *App) Run(addrWithPort string){
	fmt.Println("Server running at ", addrWithPort)
	log.Fatal(http.ListenAndServe(addrWithPort, app.Router))
}

/*
Creates the database with all its table using the sql file
A cheap migration per se (lol)
*/
func (app *App) MakeDB(user string, pw string, dbhost string, dbport string) {

	// Create the coonection string url
	// Note: multistatement and no db name has been passed in.
	conn_url := user + ":" + pw + "@tcp(" + dbhost + ":" + dbport + ")/?multiStatements=true"
	var err error
	app.DB, err = sql.Open("mysql", conn_url)
	if err != nil {
		log.Fatal(err)
	}
	// close db connection when done
	defer app.DB.Close()

	fmt.Println("Connected to mysql")

	// Read the file
	// query value is byte array
	query, err := ioutil.ReadFile("todo_db.sql")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Finished reading file")

	// execute query casted to string
	if _, err := app.DB.Exec(string(query)); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Query competed")

}

