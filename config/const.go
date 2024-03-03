package config

// envVariableName is an internal type for ENV variable names.
type EnvVariableName string

// ENV variable names
const (
	env              EnvVariableName = "ENV_NAME"
	dbType           EnvVariableName = "DB_TYPE"
	sqliteStorageDir EnvVariableName = "SQLITE_STORAGE_DIR"
	sqliteDBName     EnvVariableName = "SQLITE_DB_NAME"
	postgresUser     EnvVariableName = "POSTGRES_USER"
	postgresPassword EnvVariableName = "POSTGRES_PASSWORD"
	postgresDBName   EnvVariableName = "POSTGRES_DB"
)

const DBTypeSqlite = "sqlite"
const DBTypePostgres = "postgres"
