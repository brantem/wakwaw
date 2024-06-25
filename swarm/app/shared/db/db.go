package db

import (
	"fmt"
	"os"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func new(host, target, appID string) *sqlx.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s search_path=%s application_name=%s sslmode=disable",
		host,
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USERNAME"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DATABASE"),
		os.Getenv("PG_SCHEMA"),
		appID,
	)

	driver := otelsql.WrapDriver(
		pq.Driver{},
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithSpanOptions(otelsql.SpanOptions{
			OmitConnResetSession: true,
			OmitConnPrepare:      true,
			OmitRows:             true,
			OmitConnectorConnect: true,
		}),
	)

	logger := zerolog.New(os.Stdout)
	if os.Getenv("DEBUG") != "" {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	logger = logger.With().Str("target", target).Logger()

	attributes := []attribute.KeyValue{
		semconv.DBSystemPostgreSQL,
		attribute.String("db.target", target),
	}

	opts := []sqldblogger.Option{
		sqldblogger.WithPreparerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithQueryerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithExecerLevel(sqldblogger.LevelDebug),
	}
	_db := sqldblogger.OpenDriver(dsn, driver, zerologadapter.New(logger), opts...)
	if err := otelsql.RegisterDBStatsMetrics(_db, otelsql.WithAttributes(attributes...)); err != nil {
		log.Panic().Err(err).Msg("db.New")
	}

	db := sqlx.NewDb(_db, "postgres")
	// https://www.alexedwards.net/blog/configuring-sqldb
	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db
}

type DB struct {
	Master *sqlx.DB
	Read   *sqlx.DB
}

func New(appID string) *DB {
	masterHost := os.Getenv("PG_HOST")
	readHost := os.Getenv("PG_HOST_READ")
	if readHost == "" {
		readHost = masterHost
	}
	return &DB{
		new(masterHost, "master", appID),
		new(readHost, "read", appID),
	}
}

func (db *DB) Close() {
	db.Master.Close()
	db.Read.Close()
}
