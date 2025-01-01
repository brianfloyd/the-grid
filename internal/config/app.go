package config

import (
	"github.com/brianfloyd/the-grid/internal/db/pg"
	"github.com/brianfloyd/the-grid/internal/logger"
)

func GetPostgresConfig(config *Config) pg.PGConfiguration {
	return pg.PGConfiguration{
		Username:     config.ValueOrPanic("database.username"),
		Password:     config.ValueOrPanic("database.password"),
		Host:         config.ValueOrPanic("database.host"),
		Port:         config.ValueOrPanic("database.port"),
		DatabaseName: config.ValueOrPanic("database.name"),
		SslMode:      pg.PGSslMode(config.ValueOrPanic("database.sslmode")),
	}
}

func GetLogLevel(config *Config) logger.LogLevel {
	return logger.LogLevelFromString(config.ValueOrPanic("log.level"))
}

func GetServerAddress() string {
	return ":8080"
}
