package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/list"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/db"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/web"
	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// getLists is a handler that retrieves all rows from the list table
func (a *Application) getLists(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	lists, err := list.SelectLists(a.db)
	if err != nil {
		web.RespondError(w, r, http.StatusInternalServerError, errors.Wrap(err, "select all lists"))
	}

	if len(lists) == 0 {
		web.Respond(w, r, http.StatusNoContent, nil)
		return
	}

	web.Respond(w, r, http.StatusOK, lists)
}

// createList is a handler that inserts a new row into the list table
func (a *Application) createList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var payload list.Record

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		web.RespondError(w, r, http.StatusInternalServerError, errors.Wrap(err, "unmarshal request payload"))
		return
	}

	if payload.Name == "" {
		web.RespondError(w, r, http.StatusBadRequest, errors.New("name key is required"))
		return
	}

	l, err := list.CreateList(a.db, payload)
	if err != nil {
		if pgerr, ok := errors.Cause(err).(*pq.Error); ok {
			if string(pgerr.Code) == db.PSQLErrUniqueConstraint {
				web.RespondError(w, r, http.StatusBadRequest, errors.Wrap(err, "attempting to break unique name constraint"))
				return
			}
		}

		web.RespondError(w, r, http.StatusInternalServerError, errors.Wrap(err, "insert row into list table"))
		return
	}

	web.Respond(w, r, http.StatusCreated, l)
}

// getList is a handler that gets a single row from the list table using a given
// list id
func (a *Application) getList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	listID, err := strconv.Atoi(ps.ByName("lid"))
	if err != nil {
		web.RespondError(w, r, http.StatusInternalServerError, errors.Wrap(err, "convert list id to integer"))
		return
	}

	l, err := list.SelectList(a.db, list.FilterByID, listID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			web.RespondError(w, r, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
			return
		}

		web.RespondError(w, r, http.StatusInternalServerError, errors.Wrap(err, "select list by id"))
		return
	}

	web.Respond(w, r, http.StatusOK, l)
}

// updateList is a handler that updates a row from the list table using a given
// list id
func (a *Application) updateList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	listID, err := strconv.Atoi(ps.ByName("lid"))
	if err != nil {
		web.RespondError(w, r, http.StatusInternalServerError, errors.Wrap(err, "convert list id to integer"))
		return
	}

	var payload list.Record
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		web.RespondError(w, r, http.StatusInternalServerError, errors.Wrap(err, "unmarshal request payload"))
		return
	}

	payload.ID = listID

	if payload.Name == "" {
		web.RespondError(w, r, http.StatusBadRequest, errors.New("name key is required"))
		return
	}

	if err := list.UpdateList(a.db, payload); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			web.RespondError(w, r, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
			return
		}

		if pgerr, ok := errors.Cause(err).(*pq.Error); ok {
			if string(pgerr.Code) == db.PSQLErrUniqueConstraint {
				web.RespondError(w, r, http.StatusBadRequest, errors.Wrap(err, "attempting to break unique name constraint"))
				return
			}
		}

		web.RespondError(w, r, http.StatusInternalServerError, errors.Wrap(err, "update row in list table"))
		return
	}

	web.Respond(w, r, http.StatusOK, payload)
}

// deleteList is a handler that deletes a row from the list table using a given
// list id
func (a *Application) deleteList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	listID, err := strconv.Atoi(ps.ByName("lid"))
	if err != nil {
		web.RespondError(w, r, http.StatusInternalServerError, errors.Wrap(err, "convert list id to integer"))
		return
	}

	if err := list.DeleteList(a.db, listID); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			web.RespondError(w, r, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
			return
		}

		web.RespondError(w, r, http.StatusInternalServerError, errors.Wrap(err, "delete list by id"))
		return
	}

	web.Respond(w, r, http.StatusNoContent, nil)
}
