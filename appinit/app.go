package appinit

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	DB     *sql.DB
	Router *mux.Router
}

func (app *App) AppInit(user string, pw string, dbhost string, dbport string, dbname string) {
	cfg := mysql.Config{
		User:   user,
		Passwd: pw,
		Net:    "tcp",
		Addr:   dbhost + ":" + dbport,
		DBName: dbname,
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

	fmt.Println("Connected to ", dbname, "located at ", dbhost, ":", dbport)

	app.Router = mux.NewRouter()
}

func (app *App) AppRun(addrWithPort string) {
	fmt.Println("Server running at ", addrWithPort)
	err := http.ListenAndServe(addrWithPort, app.Router)
	if err != nil {
		log.Fatal(err)
	}

}
