package apiserver

import (
	"database/sql"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/stdlib"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/morozvol/AuthService/internal/app/store/sqlstore"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Start ...
func Start(config *Config) error {

	db, err := newDB(config, logrus.StandardLogger())

	if err != nil {
		return err
	}
	defer db.Close()

	store := sqlstore.New(db)
	srv := newServer(store, config)
	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(config *Config, logger *logrus.Logger) (*sqlx.DB, error) {

	// parse connection string
	dbConf, err := pgx.ParseConfig(config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	dbConf.Logger = logrusadapter.NewLogger(logger)

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
