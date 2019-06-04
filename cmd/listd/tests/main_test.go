package tests

import (
	"os"
	"testing"

	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/handlers"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/testdb"
	log "github.com/sirupsen/logrus"
)

// a is a reference to the main Application type. This is used for its database
// connection that it harbours inside of the type as well as the route definitions
// that are defined on the embedded handler.
var a *handlers.Application

// TestMain calls testMain and passes the returned exit code to os.Exit(). The reason
// that TestMain is basically a wrapper around testMain is because os.Exit() does not
// respect deferred functions, so this configuration allows for a deferred function.
func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

// testMain returns an integer denoting an exit code to be returned and used in
// TestMain. The exit code 0 denotes success, all other codes denote failure (1
// and 2).
func testMain(m *testing.M) int {
	dbc, err := testdb.Open()
	if err != nil {
		log.WithError(err).Info("create test database connection")
		return 1
	}
	defer dbc.Close()

	a = handlers.NewApplication(dbc)

	return m.Run()
}
