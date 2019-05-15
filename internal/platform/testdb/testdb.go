package testdb

import (
	"time"

	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/item"
	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/list"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	// DatabaseName is the name of the database that gets used during testing
	// in order to not corrupt production data just in case the tests get ran
	// with the wrong compose file.
	DatabaseName = "testdb"
)

// Truncate removes all seed data from the test database.
func Truncate(dbc *sqlx.DB) error {
	stmt := "TRUNCATE TABLE list, item;"

	if _, err := dbc.Exec(stmt); err != nil {
		return errors.Wrap(err, "truncate tables")
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
			return nil, errors.Wrap(err, "prepare seed list insertion")
		}

		row := stmt.QueryRow(lists[i].Name, lists[i].Created, lists[i].Modified)

		if err = row.Scan(&lists[i].ID); err != nil {
			if err := stmt.Close(); err != nil {
				logrus.WithField("err", err).Error("close psql statement")
			}

			return nil, errors.Wrap(err, "capture list id for seeded list")
		}

		if err := stmt.Close(); err != nil {
			logrus.WithField("err", err).Error("close psql statement")
		}
	}

	return lists, nil
}

// SeedItems handles seeding the item table in the database for integration tests.
func SeedItems(dbc *sqlx.DB, lists []list.List) ([]item.Item, error) {
	if len(lists) == 0 {
		return nil, errors.New("list data does not exist, necessary for item seeding")
	}

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
			return nil, errors.Wrap(err, "prepare seed item insertion")
		}

		row := stmt.QueryRow(items[i].ListID, items[i].Name, items[i].Quantity, items[i].Created, items[i].Modified)

		if err = row.Scan(&items[i].ID); err != nil {
			if err := stmt.Close(); err != nil {
				logrus.WithField("err", err).Error("close psql statement")
			}

			return nil, errors.Wrap(err, "capture item id for seeded item")
		}

		if err := stmt.Close(); err != nil {
			logrus.WithField("err", err).Error("close psql statement")
		}
	}

	return items, nil
}
