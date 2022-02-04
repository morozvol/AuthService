package db

import (
	"database/sql"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/morozvol/AuthService/internal/app/config"
	"github.com/sirupsen/logrus"
	"time"
)

// New return sqlx.DB
func New(config *config.Config, logger *logrus.Logger) (*sqlx.DB, error) {

	// parse connection string
	dbConf, err := pgx.ParseConfig(config.DB.GetConnactionString())
	if err != nil {
		return nil, err
	}

	dbConf.Logger = logrusadapter.NewLogger(logger)
	dbConf.Host = config.DB.Host

	// register pgx conn
	dsn := stdlib.RegisterConnConfig(dbConf)

	sql.Register("wrapper", stdlib.GetDefaultDriver())
	wdb, err := sql.Open("wrapper", dsn)
	if err != nil {
		return nil, err
	}

	db := sqlx.NewDb(wdb, "pgx")
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
