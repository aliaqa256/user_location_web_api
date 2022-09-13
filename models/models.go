package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type config struct {
	dbStr string
}

type application struct {
	config   config
	DB       *sql.DB
	errorLog *log.Logger
}

func (app *application) ConnectDB() error {
	db, err := sql.Open("postgres", app.config.dbStr)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	app.DB = db
	fmt.Println("Connected to database")
	return nil
}

func (app *application) CloseDB() error {
	fmt.Println("Closing DB")
	return app.DB.Close()
}


func (app *application) createUsersTable() {
	// create user table
	createUserTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`
	_, err := app.DB.Exec(createUserTableSQL)
	if err != nil {
		app.errorLog.Fatal(err)
	}
}


func (app *application) createUserInfoTable() {
	// create user info table
	createUserInfoTableSQL := `CREATE TABLE IF NOT EXISTS user_info (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		longitude FLOAT NOT NULL,
		latitude FLOAT NOT NULL,
		speed FLOAT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`
	_, err := app.DB.Exec(createUserInfoTableSQL)
	if err != nil {
		app.errorLog.Fatal(err)
	}
}


func (app *application)Migrate() {
	app.createUsersTable()
	app.createUserInfoTable()
}


var singularModelApp *application = nil
func NewApplication() *application {
	
	if singularModelApp == nil {
		err := godotenv.Load(".env")
		if err != nil {
			panic("Error loading .env file")
		}
		dbstr := os.Getenv("DB_STR")
		app := &application{
			config: config{
				dbStr: dbstr,
			},
		}
		singularModelApp = app
	}


	return singularModelApp
}