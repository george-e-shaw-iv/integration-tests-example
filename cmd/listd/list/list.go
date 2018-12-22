package list

import (
	"time"

	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/db"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Record is a type that contains the proper struct tags for both
// a JSON and Postgres representation of a list
type Record struct {
	ID       int       `json:"id" db:"list_id"`
	Name     string    `json:"name" db:"name"`
	Created  time.Time `json:"created" db:"created"`
	Modified time.Time `json:"modified" db:"modified"`
}

// SelectLists selects all rows from the list table
func SelectLists(dbc *sqlx.DB) ([]Record, error) {
	lists := make([]Record, 0)

	if err := dbc.Select(&lists, selectAll); err != nil {
		return nil, errors.Wrap(err, "select all rows from list table")
	}

	return lists, nil
}

// SelectList selects a single row from the list table based off given arguments and
// one of the following filters: FilterByID
func SelectList(dbc *sqlx.DB, filter string, args ...interface{}) (Record, error) {
	var (
		list Record
		stmt string
	)

	switch filter {
	case FilterByID:
		stmt = selectByID
	default:
		return Record{}, db.ErrUnknownFilter
	}

	if err := dbc.Select(&list, stmt, args...); err != nil {
		return Record{}, errors.Wrap(err, "select singular row from list table")
	}

	return list, nil
}

// CreateList inserts a new row into the list table
func CreateList(dbc *sqlx.DB, r Record) (Record, error) {
	r.Created = time.Now()
	r.Modified = time.Now()

	res, err := dbc.Exec(insert, r.Name, r.Created, r.Modified)
	if err != nil {
		return Record{}, errors.Wrap(err, "insert new list row")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Record{}, errors.Wrap(err, "get inserted row id")
	}
	r.ID = int(id)

	return r, nil
}

// UpdateList updates a row in the list table based off of a list_id. The only field
// able to be updated is the name field
func UpdateList(dbc *sqlx.DB, r Record) error {
	r.Modified = time.Now()

	if _, err := dbc.Exec(update, r.Name, r.Modified, r.ID); err != nil {
		return errors.Wrap(err, "update list row")
	}

	return nil
}

// DeleteList deletes a row in the list table based off of list_id
func DeleteList(dbc *sqlx.DB, id int) error {
	if _, err := dbc.Exec(delete, id); err != nil {
		return errors.Wrap(err, "delete list row")
	}

	return nil
}
