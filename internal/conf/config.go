package conf

import (
	"os"

	"github.com/brianfloyd/the-grid/internal/db/pg"
)

func getEnvOrPanic(key string) string {
	val, set := os.LookupEnv(key)
	if !set {
		panic("Environment variable: " + key + " was not set.")
	}
	return val
}

func LoadPgConf() pg.PGConfiguration {
	return pg.PGConfiguration{
		Username:     getEnvOrPanic("DATABASE_USERNAME"),
		Password:     getEnvOrPanic("DATABASE_PASSWORD"),
		Host:         getEnvOrPanic("DATABASE_HOST"),
		Port:         getEnvOrPanic("DATABASE_PORT"),
		DatabaseName: getEnvOrPanic("DATABASE_NAME"),
		SslMode:      pg.SslModeRequire,
	}
}

func LoadServerAddress() string {
	addr, set := os.LookupEnv("SERVER_ADDR")
	if !set {
		addr = ":8080"
	}
	return addr
}
