package config

// envVariableName is an internal type for ENV variable names.
type EnvVariableName string

// ENV variable names
const (
	env              EnvVariableName = "ENV_NAME"
	dbType           EnvVariableName = "DB_TYPE"
	sqliteStorageDir EnvVariableName = "SQLITE_STORAGE_DIR"
	sqliteDBName     EnvVariableName = "SQLITE_DB_NAME"
	pgHost           EnvVariableName = "POSTGRES_HOST"
	pgPort           EnvVariableName = "POSTGRES_PORT"
	pgUser           EnvVariableName = "POSTGRES_USER"
	pgPassword       EnvVariableName = "POSTGRES_PASSWORD"
	pgDBName         EnvVariableName = "POSTGRES_DB"
	pgSSLMode        EnvVariableName = "POSTGRES_SSL_MODE"
)

const DBTypeSqlite = "sqlite"
const DBTypePostgres = "postgres"
