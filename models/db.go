package models

import (
	"fmt"
	"path/filepath"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SQLiteConfig represents data needed to connect to SQLite
type SQLiteConfig struct {
	StorageDir  string
	Environment string
	DBName      string
}

// PostgresConfig data needed to connect to postgres
type PostgresConfig struct {
	Host     string
	Port     string
	Password string
	User     string
	SSLMode  string
	DBName   string
}

// DBConfig represents a interface for db configs
type DBConfig interface {
	ConnectToDB() (*gorm.DB, error)
}

// SQLiteConfig.ConnectToDB connects to an SQLite DB
func (cfg *SQLiteConfig) ConnectToDB() (*gorm.DB, error) {
	// TODO make it so the dir is created
	db, err := gorm.Open(
		sqlite.Open(
			filepath.Join(
				cfg.StorageDir,
				fmt.Sprintf("%s_%s.sqlite", cfg.DBName, cfg.Environment),
			),
		),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// PostgresConfig.ConnectToDB connects to a Postgres DB
func (cfg *PostgresConfig) ConnectToDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)
	// Creates connection to postgress
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil

}

// MustConnect is a Must wrapper for connector functions
func MustConnect(db *gorm.DB, err error) *gorm.DB {
	if err != nil {
		panic(err)
	}
	return db
}
