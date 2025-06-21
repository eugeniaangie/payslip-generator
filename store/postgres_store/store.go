package postgres_store

import (
	"fmt"

	"payslip-generator/util/config"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Store struct {
	logger *logrus.Logger
	Db     *sqlx.DB
}

func newDB(pgConfig config.Postgres) (*sqlx.DB, error) {
	db, err := sqlx.Open(pgConfig.Driver, pgConfig.Dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %s", err.Error())
	}

	// Optional: ping to verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return db, nil
}

func NewPostgresStore(logger *logrus.Logger, pgConfig config.Postgres) (*Store, error) {
	db, err := newDB(pgConfig)
	if err != nil {
		return nil, err
	}

	return &Store{
		logger: logger,
		Db:     db,
	}, nil
}
