package item

import (
	"time"

	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/db"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Record is a type that contains the proper struct tags for both
// a JSON and Postgres representation of an item
type Record struct {
	ID       int       `json:"id" db:"item_id"`
	ListID   int       `json:"listID" db:"list_id"`
	Name     string    `json:"name" db:"name"`
	Quantity int       `json:"quantity" db:"quantity"`
	Created  time.Time `json:"created" db:"created"`
	Modified time.Time `json:"modified" db:"modified"`
}

// SelectItems selects all appropriate rows from the item table given a list_id
func SelectItems(dbc *sqlx.DB, listID int) ([]Record, error) {
	items := make([]Record, 0)

	if err := dbc.Select(&items, selectAll, listID); err != nil {
		return nil, errors.Wrap(err, "select all rows from item table given a list_id")
	}

	return items, nil
}

// SelectItem selects a single row from the item table based off given arguments and
// one of the following filters: FilterByIDAndListID
func SelectItem(dbc *sqlx.DB, filter string, args ...interface{}) (Record, error) {
	var (
		item Record
		stmt string
	)

	switch filter {
	case FilterByIDAndListID:
		stmt = selectByIDAndListID
	default:
		return Record{}, db.ErrUnknownFilter
	}

	pStmt, err := dbc.Preparex(stmt)
	if err != nil {
		return Record{}, errors.Wrap(err, "prepare select query")
	}
	defer pStmt.Close()

	row := pStmt.QueryRowx(args...)

	if err := row.StructScan(&item); err != nil {
		return Record{}, errors.Wrap(err, "select singular row from item table")
	}

	return item, nil
}

// CreateItem inserts a new row into the item table
func CreateItem(dbc *sqlx.DB, r Record) (Record, error) {
	r.Created = time.Now()
	r.Modified = time.Now()

	stmt, err := dbc.Prepare(insert)
	if err != nil {
		return Record{}, errors.Wrap(err, "insert new item row")
	}
	defer stmt.Close()

	row := stmt.QueryRow(r.ListID, r.Name, r.Quantity, r.Created, r.Modified)

	if err = row.Scan(&r.ID); err != nil {
		return Record{}, errors.Wrap(err, "get inserted row id")
	}

	return r, nil
}

// UpdateItem updates a row in the item table based off of a item_id. The only fields
// able to be updated are the name and quantity field
func UpdateItem(dbc *sqlx.DB, r Record) error {
	r.Modified = time.Now()

	if _, err := dbc.Exec(update, r.Name, r.Quantity, r.Modified, r.ID); err != nil {
		return errors.Wrap(err, "update item row")
	}

	return nil
}

// DeleteItem deletes a row in the item table based off of item_id
func DeleteItem(dbc *sqlx.DB, id int) error {
	if _, err := dbc.Exec(del, id); err != nil {
		return errors.Wrap(err, "delete list row")
	}

	return nil
}
