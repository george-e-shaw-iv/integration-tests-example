package list

import (
	"database/sql"
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

	pStmt, err := dbc.Preparex(stmt)
	if err != nil {
		return Record{}, errors.Wrap(err, "prepare select query")
	}
	defer pStmt.Close()

	row := pStmt.QueryRowx(args...)

	if err := row.StructScan(&list); err != nil {
		return Record{}, errors.Wrap(err, "select singular row from list table")
	}

	return list, nil
}

// CreateList inserts a new row into the list table
func CreateList(dbc *sqlx.DB, r Record) (Record, error) {
	r.Created = time.Now()
	r.Modified = time.Now()

	stmt, err := dbc.Prepare(insert)
	if err != nil {
		return Record{}, errors.Wrap(err, "insert new list row")
	}
	defer stmt.Close()

	row := stmt.QueryRow(r.Name, r.Created, r.Modified)

	if err = row.Scan(&r.ID); err != nil {
		return Record{}, errors.Wrap(err, "get inserted row id")
	}

	return r, nil
}

// UpdateList updates a row in the list table based off of a list_id. The only field
// able to be updated is the name field
func UpdateList(dbc *sqlx.DB, r Record) error {
	if _, err := SelectList(dbc, FilterByID, r.ID); errors.Cause(err) == sql.ErrNoRows {
		return sql.ErrNoRows
	}

	r.Modified = time.Now()

	if _, err := dbc.Exec(update, r.Name, r.Modified, r.ID); err != nil {
		return errors.Wrap(err, "update list row")
	}

	return nil
}

// DeleteList deletes a row in the list table based off of list_id
func DeleteList(dbc *sqlx.DB, id int) error {
	if _, err := SelectList(dbc, FilterByID, id); errors.Cause(err) == sql.ErrNoRows {
		return sql.ErrNoRows
	}

	if _, err := dbc.Exec(delRelatedItems, id); err != nil && errors.Cause(err) != sql.ErrNoRows {
		return errors.Wrap(err, "deleted related items to given list_id")
	}

	if _, err := dbc.Exec(del, id); err != nil {
		return errors.Wrap(err, "delete list row")
	}

	return nil
}
