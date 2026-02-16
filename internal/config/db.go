package config

type DBConfig struct {
	SQLitePath string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		SQLitePath: getEnvDefault("SQLITE_PATH", SQLitePath),
	}
}