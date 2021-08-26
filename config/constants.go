package config

const (
	EnvironmentDefault     = "default"
	EnvironmentProduction  = "production"
	EnvironmentTest        = "test"
	EnvironmentDevelopment = "development"
	EnvironmentLocal       = "local"
	EnvironmentDocker      = "docker"
	EnvironmentCiMongo     = "ci_mongo"
	EnvironmentCiSQLites3  = "ci_sqlite3"

	AppSettingsEnvName = "ENVIRONMENT"

	DbTypeMongo   = "mongo"
	DbTypeSQLite3 = "sqlite3"

	NotifierTypeTelegram = "telegram"
)
