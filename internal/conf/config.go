package conf

import (
	"fmt"
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

type Domain string

const (
	TEST Domain = "test"
	PROD Domain = "prod"
)

func getDomain() Domain {
	domain := (Domain)(getEnvOrPanic("DOMAIN"))
	if domain != TEST && domain != PROD {
		panic(fmt.Sprintf("Unknown domain: %s", domain))
	}
	return domain
}

func isTest() bool {
	return getDomain() == TEST
}

func LoadPgConf() pg.PGConfiguration {
	getSslMode := func() pg.PGSslMode {
		if isTest() {
			return pg.SslModeDisable
		} else {
			return pg.SslModeRequire
		}
	}

	return pg.PGConfiguration{
		Username:     getEnvOrPanic("DATABASE_USERNAME"),
		Password:     getEnvOrPanic("DATABASE_PASSWORD"),
		Host:         getEnvOrPanic("DATABASE_HOST"),
		Port:         getEnvOrPanic("DATABASE_PORT"),
		DatabaseName: getEnvOrPanic("DATABASE_NAME"),
		SslMode:      getSslMode(),
	}
}

func LoadServerAddress() string {
	addr, set := os.LookupEnv("SERVER_ADDR")
	if !set {
		addr = ":8080"
	}
	return addr
}
