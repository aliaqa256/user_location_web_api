package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	dbStr string

}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
	db       *sql.DB
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Starting Back end server in %s mode on port %d\n", app.config.env, app.config.port)

	return srv.ListenAndServe()
}

func main() {
	var cfg config
	cfg.dbStr="host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"

	flag.IntVar(&cfg.port, "port", 4001, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production|maintenance}")

	flag.Parse()



	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)


	// connect to the database
	db, err := sql.Open("postgres", cfg.dbStr)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()


	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
		db:       db,
	}
	// create db models
	app.Migrate()

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}


}
