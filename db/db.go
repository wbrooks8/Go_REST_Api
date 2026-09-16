package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the shared database handle used by the model layer.
// *sql.DB manages a pool of connections; it is not a single open connection.
var DB *sql.DB

// InitDB opens the SQLite database file and creates the tables needed by the
// application. Call this once during application startup before serving routes.
func InitDB() {
	var err error
	// sql.Open configures a handle. The first actual database work may happen
	// later, which is why Exec errors must still be checked below.
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("Could not connect to database")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	// Create tables after opening the database so models can use them immediately.
	createTables()
}

func createTables() {
	// Exec runs a SQL statement that does not return rows. IF NOT EXISTS makes
	// startup safe to repeat without trying to recreate existing tables.
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)
	`

	_, err := DB.Exec(createUsersTable)

	if err != nil {
		panic("Could not create users table")
	}

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		userID INTEGER, 
		FOREIGN KEY (userID) REFERENCES users(id)
	)
	`

	_, err = DB.Exec(createEventsTable)

	if err != nil {
		panic("Could not create events table")
	}

	// This join table connects users and events. One row means that one user
	// registered for one event. The existing name is misspelled but retained so
	// it remains compatible with the model queries and existing database files.
	createRegistrationsTable := `
	CREATE TABLE IF NOT EXISTS regestrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_id INTEGER,
		user_id INTEGER,
		FOREIGN KEY(event_id) REFERENCES events(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`

	_, err = DB.Exec(createRegistrationsTable)
	if err != nil {
		panic("Could not create regestrations table")
	}
}
