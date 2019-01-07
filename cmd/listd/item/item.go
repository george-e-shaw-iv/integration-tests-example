package item

import (
	"database/sql"
	"time"

	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/list"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Item is a type that contains the proper struct tags for both
// a JSON and Postgres representation of an item.
type Item struct {
	ID       int       `json:"id" db:"item_id"`
	ListID   int       `json:"listID" db:"list_id"`
	Name     string    `json:"name" db:"name"`
	Quantity int       `json:"quantity" db:"quantity"`
	Created  time.Time `json:"created" db:"created"`
	Modified time.Time `json:"modified" db:"modified"`
}

// SelectItems selects all appropriate rows from the item table given a list_id.
func SelectItems(dbc *sqlx.DB, listID int) ([]Item, error) {
	if _, err := list.SelectList(dbc, listID); errors.Cause(err) == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}

	items := make([]Item, 0)

	if err := dbc.Select(&items, selectAll, listID); err != nil {
		return nil, errors.Wrap(err, "select all rows from item table given a list_id")
	}

	return items, nil
}

// SelectItem selects a single row from the item table based off given list_id and
// item_id.
func SelectItem(dbc *sqlx.DB, iid, lid int) (Item, error) {
	var i Item
	stmt := selectByIDAndListID

	pStmt, err := dbc.Preparex(stmt)
	if err != nil {
		return Item{}, errors.Wrap(err, "prepare select query")
	}

	defer func() {
		if err := pStmt.Close(); err != nil {
			logrus.WithError(errors.Wrap(err, "close psql statement")).Info("select item")
		}
	}()

	row := pStmt.QueryRowx(iid, lid)

	if err := row.StructScan(&i); err != nil {
		return Item{}, errors.Wrap(err, "select singular row from item table")
	}

	return i, nil
}

// CreateItem inserts a new row into the item table.
func CreateItem(dbc *sqlx.DB, r Item) (Item, error) {
	r.Created = time.Now()
	r.Modified = time.Now()

	if _, err := list.SelectList(dbc, r.ListID); errors.Cause(err) == sql.ErrNoRows {
		return Item{}, sql.ErrNoRows
	}

	stmt, err := dbc.Prepare(insert)
	if err != nil {
		return Item{}, errors.Wrap(err, "insert new item row")
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			logrus.WithError(errors.Wrap(err, "close psql statement")).Info("create item")
		}
	}()

	row := stmt.QueryRow(r.ListID, r.Name, r.Quantity, r.Created, r.Modified)

	if err = row.Scan(&r.ID); err != nil {
		return Item{}, errors.Wrap(err, "get inserted row id")
	}

	return r, nil
}

// UpdateItem updates a row in the item table based off of item_id and list_id. The only fields
// able to be updated are the name and quantity field.
func UpdateItem(dbc *sqlx.DB, r Item) error {
	if _, err := SelectItem(dbc, r.ID, r.ListID); errors.Cause(err) == sql.ErrNoRows {
		return sql.ErrNoRows
	}

	r.Modified = time.Now()

	if _, err := dbc.Exec(update, r.Name, r.Quantity, r.Modified, r.ID, r.ListID); err != nil {
		return errors.Wrap(err, "update item row")
	}

	return nil
}

// DeleteItem deletes a row in the item table based off of item_id.
func DeleteItem(dbc *sqlx.DB, itemID, listID int) error {
	if _, err := SelectItem(dbc, itemID, listID); errors.Cause(err) == sql.ErrNoRows {
		return sql.ErrNoRows
	}

	if _, err := dbc.Exec(del, itemID); err != nil {
		return errors.Wrap(err, "delete list row")
	}

	return nil
}
