package main


// create user table
func (app *application) createUsersTable() {
	// create user table
	createUserTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`
	_, err := app.db.Exec(createUserTableSQL)
	if err != nil {
		app.errorLog.Fatal(err)
	}
}

// create user info table
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
	_, err := app.db.Exec(createUserInfoTableSQL)
	if err != nil {
		app.errorLog.Fatal(err)
	}
}

func (app *application)Migrate() {
	app.createUsersTable()
	app.createUserInfoTable()
}