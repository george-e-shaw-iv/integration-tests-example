package list

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// List is a type that contains the proper struct tags for both
// a JSON and Postgres representation of a list.
type List struct {
	ID       int       `json:"id" db:"list_id"`
	Name     string    `json:"name" db:"name"`
	Created  time.Time `json:"created" db:"created"`
	Modified time.Time `json:"modified" db:"modified"`
}

// SelectLists selects all rows from the list table.
func SelectLists(dbc *sqlx.DB) ([]List, error) {
	lists := make([]List, 0)

	if err := dbc.Select(&lists, selectAll); err != nil {
		return nil, errors.Wrap(err, "select all rows from list table")
	}

	return lists, nil
}

// SelectList selects a single row from the list table based off of a given list_id.
func SelectList(dbc *sqlx.DB, id int) (List, error) {
	var list List
	stmt := selectByID

	pStmt, err := dbc.Preparex(stmt)
	if err != nil {
		return List{}, errors.Wrap(err, "prepare select query")
	}

	defer func() {
		if err := pStmt.Close(); err != nil {
			logrus.WithError(errors.Wrap(err, "close psql statement")).Info("select list")
		}
	}()

	row := pStmt.QueryRowx(id)

	if err := row.StructScan(&list); err != nil {
		return List{}, errors.Wrap(err, "select singular row from list table")
	}

	return list, nil
}

// CreateList inserts a new row into the list table.
func CreateList(dbc *sqlx.DB, r List) (List, error) {
	r.Created = time.Now()
	r.Modified = time.Now()

	stmt, err := dbc.Prepare(insert)
	if err != nil {
		return List{}, errors.Wrap(err, "insert new list row")
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			logrus.WithError(errors.Wrap(err, "close psql statement")).Info("create list")
		}
	}()

	row := stmt.QueryRow(r.Name, r.Created, r.Modified)

	if err = row.Scan(&r.ID); err != nil {
		return List{}, errors.Wrap(err, "get inserted row id")
	}

	return r, nil
}

// UpdateList updates a row in the list table based off of a list_id. The only field
// able to be updated is the name field.
func UpdateList(dbc *sqlx.DB, r List) error {
	if _, err := SelectList(dbc, r.ID); errors.Cause(err) == sql.ErrNoRows {
		return sql.ErrNoRows
	}

	r.Modified = time.Now()

	if _, err := dbc.Exec(update, r.Name, r.Modified, r.ID); err != nil {
		return errors.Wrap(err, "update list row")
	}

	return nil
}

// DeleteList deletes a row in the list table based off of list_id.
func DeleteList(dbc *sqlx.DB, id int) error {
	if _, err := SelectList(dbc, id); errors.Cause(err) == sql.ErrNoRows {
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
