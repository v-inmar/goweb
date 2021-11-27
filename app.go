package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	fmt.Println("Successfully connected to database ", dbname, " located at ", dbhost, ":", dbport)
	app.Router = mux.NewRouter()
}

func (app *App) Run(addrWithPort string){

	// https://github.com/gorilla/mux#graceful-shutdown
	var wait time.Duration
    flag.DurationVar(&wait, "graceful-timeout", time.Second * 15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
    flag.Parse()

	srvr := &http.Server{
		Addr: addrWithPort,
		WriteTimeout: time.Second * 15,
		ReadTimeout: time.Second * 15,
		IdleTimeout: time.Second * 60,
		Handler: app.Router,
	}

	
	go func ()  {
		fmt.Println("Server running at ", addrWithPort)
		if err := srvr.ListenAndServe(); err != nil{
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srvr.Shutdown(ctx)
	log.Println("Shutting down")
	os.Exit(0)
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

