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

// Seed handles seeding all necessary tables in the database in order to carry
// out integration testing.
func Seed(dbc *sqlx.DB) ([]list.List, []item.Item, error) {
	if err := Truncate(dbc); err != nil {
		return nil, nil, errors.Wrap(err, "truncate before seeding")
	}

	now := time.Now().UTC().Truncate(time.Microsecond)

	lists, err := seedLists(dbc, now)
	if err != nil {
		return nil, nil, errors.Wrap(err, "seed list data")
	}

	items, err := seedItems(dbc, now, lists)
	if err != nil {
		return nil, nil, errors.Wrap(err, "seed item data")
	}

	return lists, items, nil
}

// Truncate removes all seed data from the test database.
func Truncate(dbc *sqlx.DB) error {
	stmt := "TRUNCATE TABLE list, item;"

	if _, err := dbc.Exec(stmt); err != nil {
		return errors.Wrap(err, "truncate tables")
	}

	return nil
}

// seedLists handles seeding the list table in the database for integration tests.
func seedLists(dbc *sqlx.DB, t time.Time) ([]list.List, error) {
	lists := []list.List{
		{
			Name:     "Grocery",
			Created:  t,
			Modified: t,
		},
		{
			Name:     "Todo",
			Created:  t,
			Modified: t,
		},
		{
			Name:     "Employees",
			Created:  t,
			Modified: t,
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

// seedItems handles seeding the item table in the database for integration tests.
func seedItems(dbc *sqlx.DB, t time.Time, lists []list.List) ([]item.Item, error) {
	if len(lists) == 0 {
		return nil, errors.New("list data does not exist, necessary for item seeding")
	}

	items := []item.Item{
		{
			ListID:   lists[0].ID, // Grocery List
			Name:     "Chocolate Milk",
			Quantity: 1,
			Created:  t,
			Modified: t,
		},
		{
			ListID:   lists[0].ID, // Grocery List
			Name:     "Mac and Cheese",
			Quantity: 2,
			Created:  t,
			Modified: t,
		},
		{
			ListID:   lists[1].ID, // Todo List
			Name:     "Write Integration Tests",
			Quantity: 1,
			Created:  t,
			Modified: t,
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
