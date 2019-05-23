package testdb

import (
	"testing"
	"time"

	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/item"
	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/list"
	"github.com/jmoiron/sqlx"
)

const (
	// DatabaseName is the name of the database that gets used during testing
	// in order to not corrupt production data just in case the tests get ran
	// with the wrong compose file.
	DatabaseName = "testdb"
)

// Truncate removes all seed data from the test database.
func Truncate(t *testing.T, dbc *sqlx.DB) {
	t.Helper()

	stmt := "TRUNCATE TABLE list, item;"

	if _, err := dbc.Exec(stmt); err != nil {
		t.Fatalf("error truncating database tables: %v", err)
	}
}

// SeedLists handles seeding the list table in the database for integration tests.
func SeedLists(t *testing.T, dbc *sqlx.DB) []list.List {
	t.Helper()

	now := time.Now().Truncate(time.Microsecond)

	lists := []list.List{
		{
			Name:     "Grocery",
			Created:  now,
			Modified: now,
		},
		{
			Name:     "Todo",
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
			t.Errorf("error preparing seed list insertion: %v", err)
		}

		row := stmt.QueryRow(lists[i].Name, lists[i].Created, lists[i].Modified)

		if err = row.Scan(&lists[i].ID); err != nil {
			if err := stmt.Close(); err != nil {
				t.Errorf("error closing psql statement: %v", err)
			}

			t.Errorf("error capturing list id for seeded list: %v", err)
		}

		if err := stmt.Close(); err != nil {
			t.Errorf("error closing psql statement: %v", err)
		}
	}

	return lists
}

// SeedItems handles seeding the item table in the database for integration tests.
func SeedItems(t *testing.T, dbc *sqlx.DB, lists []list.List) []item.Item {
	t.Helper()

	now := time.Now().Truncate(time.Microsecond)

	items := []item.Item{
		{
			ListID:   lists[0].ID, // Grocery List
			Name:     "Chocolate Milk",
			Quantity: 1,
			Created:  now,
			Modified: now,
		},
		{
			ListID:   lists[0].ID, // Grocery List
			Name:     "Mac and Cheese",
			Quantity: 2,
			Created:  now,
			Modified: now,
		},
		{
			ListID:   lists[1].ID, // Todo List
			Name:     "Write Integration Tests",
			Quantity: 1,
			Created:  now,
			Modified: now,
		},
	}

	for i := range items {
		stmt, err := dbc.Prepare("INSERT INTO item (list_id, name, quantity, created, modified) VALUES ($1, $2, $3, $4, $5) RETURNING item_id;")
		if err != nil {
			t.Errorf("error preparing seed item insertion: %v", err)
		}

		row := stmt.QueryRow(items[i].ListID, items[i].Name, items[i].Quantity, items[i].Created, items[i].Modified)

		if err = row.Scan(&items[i].ID); err != nil {
			if err := stmt.Close(); err != nil {
				t.Errorf("error closing psql statement: %v", err)
			}

			t.Errorf("error capturing list id for seeded item: %v", err)
		}

		if err := stmt.Close(); err != nil {
			t.Errorf("error closing psql statement: %v", err)
		}
	}

	return items
}
