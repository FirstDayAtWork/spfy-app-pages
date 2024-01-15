package repository

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const SQLitePath string = "local_storage/db.sqlite"

func Connect(DBPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(DBPath), &gorm.Config{})
	if err != nil {
		// TODO Log this!
		return nil, err
	}
	return db, nil
}
