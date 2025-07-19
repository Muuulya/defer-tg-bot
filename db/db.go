package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Create tables
	schema := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY,
        telegram_id INTEGER UNIQUE,
        state TEXT
    );
    CREATE TABLE IF NOT EXISTS channels (
        id INTEGER PRIMARY KEY,
        telegram_id INTEGER,
        user_id INTEGER,
        name TEXT
    );
    CREATE TABLE IF NOT EXISTS scheduled_posts (
        id INTEGER PRIMARY KEY,
        user_id INTEGER,
        channel_id INTEGER,
        message_ids TEXT,
        scheduled_time DATETIME
    );`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}
