package db

import (
	"fmt"
	"time"

	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/configuration"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrUnknownFilter is an error that denotes an unknown filter was given to a function
	// that performs a database query that uses filters
	ErrUnknownFilter = errors.New("unknown filter given")
)

var (
	// PSQLErrUniqueConstraint holds the error code that denotes a unique constraint is
	// attempting to be violated
	PSQLErrUniqueConstraint = "23505"
)

// NewConnection returns a new database connection with the schema applied
func NewConnection(cfg *configuration.Config) (*sqlx.DB, error) {
	conn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBHost)

	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return nil, errors.Wrap(err, "connect to postgres database")
	}

	log.Info("verifying postgres connection...")
	if err := db.Ping(); err != nil {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		for range ticker.C {
			if err := db.Ping(); err == nil {
				break
			}
		}
	}
	log.Info("connected to postgres database")

	if _, err = db.Exec(schema); err != nil {
		return nil, errors.Wrap(err, "apply database schema")
	}

	return db, nil
}