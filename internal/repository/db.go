package repository

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS materials (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    unit TEXT NOT NULL,
    current_stock REAL NOT NULL,
    min_stock REAL NOT NULL
);`

	_, err = db.Exec(query)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
