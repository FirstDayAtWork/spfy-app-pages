package config

import (
	"errors"
	"log"
	"os"

	"github.com/FirstDayAtWork/mustracker/models"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	DBConfig
}

func readSQLiteConfig() (DBConfig, error) {
	if os.Getenv(string(sqliteStorageDir)) == models.EmptyString {
		return nil, errors.New("sqlite storage dir env variable is not set")
	}
	res := &SQLiteConfig{}
	res.StorageDir = os.Getenv(string(sqliteStorageDir))

	if os.Getenv(string(sqliteDBName)) == models.EmptyString {
		return nil, errors.New("sqlite db name is not set")
	}
	res.DBName = os.Getenv(string(sqliteDBName))

	if os.Getenv(string(env)) == models.EmptyString {
		return nil, errors.New("env name is not set")
	}
	res.Environment = os.Getenv(string(env))

	return res, nil
}

func readPostgresConfig() (DBConfig, error) {
	if os.Getenv(string(pgHost)) == models.EmptyString {
		return nil, errors.New("postgres host env var is not set")
	}
	res := &PostgresConfig{}
	res.Host = os.Getenv(string(pgHost))

	if os.Getenv(string(pgPort)) == models.EmptyString {
		return nil, errors.New("postgres host env var is not set")
	}
	res.Port = os.Getenv(string(pgPort))

	if os.Getenv(string(pgUser)) == models.EmptyString {
		return nil, errors.New("postgres host env var is not set")
	}
	res.User = os.Getenv(string(pgUser))

	if os.Getenv(string(pgPassword)) == models.EmptyString {
		return nil, errors.New("postgres host env var is not set")
	}
	res.Password = os.Getenv(string(pgPassword))

	if os.Getenv(string(pgDBName)) == models.EmptyString {
		return nil, errors.New("postgres host env var is not set")
	}
	res.DBName = os.Getenv(string(pgDBName))

	if os.Getenv(string(pgSSLMode)) == models.EmptyString {
		return nil, errors.New("postgres host env var is not set")
	}
	res.SSLMode = os.Getenv(string(pgSSLMode))

	return res, nil
}

func readDBConfig() (DBConfig, error) {
	switch os.Getenv(string(dbType)) {
	case DBTypeSqlite:
		log.Println("db to use is sqlite, configuring")
		return readSQLiteConfig()
	case DBTypePostgres:
		log.Println("db to use is postgres, configuring")
		return readPostgresConfig()
	default:
		log.Printf("unexpected db type: %s\n", string(os.Getenv(string(dbType))))
		return nil, errors.New("db type is not postgres or sqlite")
	}

}

func ReadConfig() (*AppConfig, error) {
	godotenv.Load()
	log.Println("read env variables")

	dbCfg, err := readDBConfig()
	if err != nil {
		log.Printf("error reading db config: %s\n", err)
		return nil, err
	}
	log.Println("successfully configured db")

	return &AppConfig{
		dbCfg,
	}, nil
}
