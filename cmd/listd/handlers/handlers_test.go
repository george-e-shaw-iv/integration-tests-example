package handlers

import (
	"os"
	"testing"

	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/configuration"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/db"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/testdb"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// testSuite is a struct type that contains necessary fields to carry out
// tasks to fully test the handlers package along with it's integrations.
type testSuite struct {
	a *Application
}

// ts is the global variable that is of type testSuite which helps test the
// entirety of the handlers package and it's integrations.
var ts testSuite

// TestMain handles the setup of the testSuite, runs all of the unit tests within
// the handlers package, and cleans up afterward.
func TestMain(m *testing.M) {
	var err error
	var dbc *sqlx.DB

	exitCode := 1

	defer func() {
		if err != nil {
			log.WithError(err).Info("error in handlers TestMain")
		}

		if dbc != nil {
			if err = dbc.Close(); err != nil {
				log.WithError(err).Info("close test database connection")
			}
		}

		os.Exit(exitCode)
	}()

	if dbc, err = db.NewConnection(&configuration.Config{
		DBUser: configuration.DefaultDBUser,
		DBPass: configuration.DefaultDBPass,
		DBName: testdb.DatabaseName,
		DBHost: configuration.DefaultDBHost,
		DBPort: configuration.DefaultDBPort,
	}); err != nil {
		err = errors.Wrap(err, "create test database connection")
		return
	}

	ts.a = NewApplication(dbc, &configuration.Config{})

	exitCode = m.Run()
}
