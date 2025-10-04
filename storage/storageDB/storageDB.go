package storageDB

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/muuulya/defer-tg-bot/storage"
)

const (
	dbPath = "storage.db"
)

type StorageDB struct {
	db *sql.DB
}

func NewStorageDB() (*StorageDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return nil, storage.ErrorCreateStorage
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if _, err = db.Exec(initSQL); err != nil {
		log.Println(err)
		return nil, storage.ErrorInitStorage
	}

	if _, err = db.Exec(createTablsSQL); err != nil {
		log.Println(err)
		return nil, storage.ErrorInitStorage
	}

	storage := StorageDB{db: db}
	return &storage, nil
}

func (s *StorageDB) Close() { s.db.Close() }

type scanner interface {
	Scan(dest ...any) error
}
