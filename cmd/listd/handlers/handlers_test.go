package handlers

import (
	"os"
	"testing"

	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/configuration"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/db"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/testdb"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// testSuite is a struct type that contains necessary fields to carry out
// tasks to fully test the handlers package along with it's integrations
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
// entirety of the handlers package and it's integrations
var ts testSuite

// TestMain handles the setup of the testSuite, runs all of the unit tests within
// the handlers package, and cleans up afterward
func TestMain(m *testing.M) {
	var mainErr error
	defer func() {
		if mainErr != nil {
			log.WithFields(log.Fields{
				"error": mainErr,
			}).Error("error in handlers TestMain")

			os.Exit(1)
		}
	}()

	dbc, err := db.NewConnection(&configuration.Config{
		DBUser: configuration.DefaultDBUser,
		DBPass: configuration.DefaultDBPass,
		DBName: testdb.DatabaseName,
		DBHost: configuration.DefaultDBHost,
		DBPort: configuration.DefaultDBPort,
	})
	if err != nil {
		mainErr = errors.Wrap(err, "create test database connection")
		return
	}

	ts.a = NewApplication(dbc, &configuration.Config{})

	// Initial test seeding
	if err := testdb.Seed(ts.a.db); err != nil {
		mainErr = errors.Wrap(err, "seeding test database")
		return
	}

	code := m.Run()

	// Clean-up
	if err := dbc.Close(); err != nil {
		log.Printf("error closing database connection: %v", err)
	}

	// m.Run() and os.Exit have to be separated between the clean up code
	// because os.Exit does not respect deferred statements/functions.
	os.Exit(code)
}
