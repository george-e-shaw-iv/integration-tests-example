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

// reseedDatabase is a function attached to the testSuite type that attempts
// to reseed the database back to its original testing state.
func (ts *testSuite) reseedDatabase(t *testing.T) {
	if err := testdb.Seed(ts.a.db); err != nil {
		t.Errorf("error encountered while seeding database: %v", err)
	}
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

	// Initial seeding of the test database using test values defined within
	// the testdb package. The testdb.Seed function also truncates all tables
	// before seeding them.
	if err = testdb.Seed(ts.a.db); err != nil {
		err = errors.Wrap(err, "seeding test database")
		return
	}

	exitCode = m.Run()
}
