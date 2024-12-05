package pg

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConnection interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
}

func New(db DbConnection) *Queries {
	return &Queries{
		db:      db,
		workout: WorkoutQueries{},
		user:    UserQueries{},
	}
}

type Queries struct {
	db      DbConnection
	workout WorkoutQueries
	user    UserQueries
}

type PGSslMode string

var (
	SslModeDisable    PGSslMode = "disable"
	SslModeAllow      PGSslMode = "allow"
	SslModePrefer     PGSslMode = "prefer"
	SslModeRequire    PGSslMode = "require"
	SslModeVerifyCa   PGSslMode = "verify-ca"
	SslModeVerifyFull PGSslMode = "verify-full"
)

type PGConfiguration struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DatabaseName string
	SslMode      PGSslMode
}

func NewPostgreSQL(conf PGConfiguration) (*pgxpool.Pool, error) {
	url := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(conf.Username, conf.Password),
		Host:   fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Path:   conf.DatabaseName,
	}

	q := url.Query()
	q.Add("sslmode", string(conf.SslMode))
	url.RawQuery = q.Encode()

	pool, err := pgxpool.New(context.TODO(), url.String())
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.TODO()); err != nil {
		return nil, err
	}

	return pool, nil
}
