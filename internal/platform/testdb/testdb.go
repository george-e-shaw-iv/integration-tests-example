package testdb

import (
	"time"

	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/item"
	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/list"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/db"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// These constants define the database connection variables for the test database
// which are used in order to not corrupt production data just in case the tests
// get ran with the wrong compose file.
const (
	// databaseUser is the user for the test database.
	databaseUser = "root"

	// databasePass is the password of the user for the test database.
	databasePass = "root"

	// databaseName is the name of the test database.
	databaseName = "testdb"

	// databaseHost is the host name of the test database.
	databaseHost = "db"

	// databasePort is the port that the test database is listening on.
	databasePort = 5432
)

// Open returns a new database connection for the test database.
func Open() (*sqlx.DB, error) {
	return db.NewConnection(db.Config{
		User: databaseUser,
		Pass: databasePass,
		Name: databaseName,
		Host: databaseHost,
		Port: databasePort,
	})
}

// Truncate removes all seed data from the test database.
func Truncate(dbc *sqlx.DB) error {
	stmt := "TRUNCATE TABLE list, item;"

	if _, err := dbc.Exec(stmt); err != nil {
		return errors.Wrap(err, "truncate test database tables")
	}

	return nil
}

// SeedLists handles seeding the list table in the database for integration tests.
func SeedLists(dbc *sqlx.DB) ([]list.List, error) {
	now := time.Now().Truncate(time.Microsecond)

	lists := []list.List{
		{
			Name:     "Grocery",
			Created:  now,
			Modified: now,
		},
		{
			Name:     "To-do",
			Created:  now,
			Modified: now,
		},
		{
			Name:     "Employees",
			Created:  now,
			Modified: now,
		},
	}

	for i := range lists {
		stmt, err := dbc.Prepare("INSERT INTO list (name, created, modified) VALUES ($1, $2, $3) RETURNING list_id;")
		if err != nil {
			return nil, errors.Wrap(err, "prepare list insertion")
		}

		row := stmt.QueryRow(lists[i].Name, lists[i].Created, lists[i].Modified)

		if err = row.Scan(&lists[i].ID); err != nil {
			if err := stmt.Close(); err != nil {
				return nil, errors.Wrap(err, "close psql statement")
			}

			return nil, errors.Wrap(err, "capture list id")
		}

		if err := stmt.Close(); err != nil {
			return nil, errors.Wrap(err, "close psql statement")
		}
	}

	return lists, nil
}

// SeedItems handles seeding the item table in the database for integration tests.
func SeedItems(dbc *sqlx.DB, lists []list.List) ([]item.Item, error) {
	now := time.Now().Truncate(time.Microsecond)

	items := []item.Item{
		{
			ListID:   lists[0].ID, // Grocery
			Name:     "Chocolate Milk",
			Quantity: 1,
			Created:  now,
			Modified: now,
		},
		{
			ListID:   lists[0].ID, // Grocery
			Name:     "Mac and Cheese",
			Quantity: 2,
			Created:  now,
			Modified: now,
		},
		{
			ListID:   lists[1].ID, // To-do
			Name:     "Write Integration Tests",
			Quantity: 1,
			Created:  now,
			Modified: now,
		},
	}

	for i := range items {
		stmt, err := dbc.Prepare("INSERT INTO item (list_id, name, quantity, created, modified) VALUES ($1, $2, $3, $4, $5) RETURNING item_id;")
		if err != nil {
			return nil, errors.Wrap(err, "prepare item insertion")
		}

		row := stmt.QueryRow(items[i].ListID, items[i].Name, items[i].Quantity, items[i].Created, items[i].Modified)

		if err = row.Scan(&items[i].ID); err != nil {
			if err := stmt.Close(); err != nil {
				return nil, errors.Wrap(err, "close psql statement")
			}

			return nil, errors.Wrap(err, "capture list id")
		}

		if err := stmt.Close(); err != nil {
			return nil, errors.Wrap(err, "close psql statement")
		}
	}

	return items, nil
}
