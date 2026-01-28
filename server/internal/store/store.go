package store

import (
	"database/sql"
	"log"
)

type Store struct {
	*sql.DB
}

func New() (*Store, error) {
	// Creates greeter.db if it doesn't exist
	db, err := sql.Open("sqlite3", "./greeter.db")
	if err != nil {
		return nil, err
	}

	// Create table if it doesn't exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		status INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return &Store{db}, nil
}
