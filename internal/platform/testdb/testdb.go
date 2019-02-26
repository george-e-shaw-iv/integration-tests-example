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

var (
	// SeedLists contain seed data for the list table.
	SeedLists []list.List

	// SeedItems contain seed data for the item table.
	SeedItems []item.Item
)

// Seed handles seeding all necessary tables in the database in order to carry
// out integration testing.
func Seed(dbc *sqlx.DB) error {
	if err := Truncate(dbc); err != nil {
		return errors.Wrap(err, "truncate before seeding")
	}

	now := time.Now().UTC().Truncate(time.Second)

	if err := seedLists(dbc, now); err != nil {
		return errors.Wrap(err, "seed list data")
	}

	if err := seedItems(dbc, now); err != nil {
		return errors.Wrap(err, "seed item data")
	}

	return nil
}

// Truncate removes all seed data from the test database.
func Truncate(dbc *sqlx.DB) error {
	stmt := "TRUNCATE TABLE list CASCADE;"

	if _, err := dbc.Exec(stmt); err != nil {
		return errors.Wrap(err, "truncate tables")
	}

	return nil
}

// seedItems handles seeding the list table in the database for integration tests.
func seedLists(dbc *sqlx.DB, t time.Time) error {
	SeedLists = []list.List{
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

	for k, l := range SeedLists {
		stmt, err := dbc.Prepare("INSERT INTO list (name, created, modified) VALUES ($1, $2, $3) RETURNING list_id;")
		if err != nil {
			return errors.Wrap(err, "prepare seed list insertion")
		}

		row := stmt.QueryRow(l.Name, l.Created, l.Modified)

		if err = row.Scan(&SeedLists[k].ID); err != nil {
			if err := stmt.Close(); err != nil {
				logrus.WithField("err", err).Error("close psql statement")
			}

			return errors.Wrap(err, "capture list id for seeded list")
		}

		if err := stmt.Close(); err != nil {
			logrus.WithField("err", err).Error("close psql statement")
		}
	}

	return nil
}

// seedItems handles seeding the item table in the database for integration tests.
func seedItems(dbc *sqlx.DB, t time.Time) error {
	if len(SeedLists) == 0 {
		return errors.New("list data does not exist, necessary for item seeding")
	}

	SeedItems = []item.Item{
		{
			ListID:   SeedLists[0].ID, // Grocery List
			Name:     "Chocolate Milk",
			Quantity: 1,
			Created:  t,
			Modified: t,
		},
		{
			ListID:   SeedLists[0].ID, // Grocery List
			Name:     "Mac and Cheese",
			Quantity: 2,
			Created:  t,
			Modified: t,
		},
		{
			ListID:   SeedLists[1].ID, // Todo List
			Name:     "Write Integration Tests",
			Quantity: 1,
			Created:  t,
			Modified: t,
		},
	}

	for k, i := range SeedItems {
		stmt, err := dbc.Prepare("INSERT INTO item (list_id, name, quantity, created, modified) VALUES ($1, $2, $3, $4, $5) RETURNING item_id;")
		if err != nil {
			return errors.Wrap(err, "prepare seed item insertion")
		}

		row := stmt.QueryRow(i.ListID, i.Name, i.Quantity, i.Created, i.Modified)

		if err = row.Scan(&SeedItems[k].ID); err != nil {
			if err := stmt.Close(); err != nil {
				logrus.WithField("err", err).Error("close psql statement")
			}

			return errors.Wrap(err, "capture item id for seeded item")
		}

		if err := stmt.Close(); err != nil {
			logrus.WithField("err", err).Error("close psql statement")
		}
	}

	return nil
}
