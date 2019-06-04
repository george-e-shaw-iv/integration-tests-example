package handlers

import (
	"net/http"

	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

// Application is the struct that contains the server handler as well as
// any references to services that the application needs.
type Application struct {
	http.Handler
	DB *sqlx.DB
}

// NewApplication returns a new pointer to Application with route definitions
// initiated.
func NewApplication(db *sqlx.DB) *Application {
	a := Application{
		DB: db,
	}

	r := httprouter.New()

	probeHandler := func(w http.ResponseWriter, r *http.Request) {
		if err := a.DB.Ping(); err == nil {

			// Ping by itself is un-reliable, the connections are cached. This
			// ensures that the database is still running by executing a harmless
			// dummy query against it.
			if _, err = a.DB.Exec("SELECT true"); err == nil {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		w.WriteHeader(http.StatusInternalServerError)
	}

	// Kubernetes Probes
	r.HandlerFunc(http.MethodGet, "/ready", probeHandler)
	r.HandlerFunc(http.MethodGet, "/healthy", probeHandler)

	// List Routes
	r.HandlerFunc(http.MethodGet, "/list", a.getLists)
	r.HandlerFunc(http.MethodPost, "/list", a.createList)
	r.HandlerFunc(http.MethodGet, "/list/:lid", a.getList)
	r.HandlerFunc(http.MethodPut, "/list/:lid", a.updateList)
	r.HandlerFunc(http.MethodDelete, "/list/:lid", a.deleteList)

	// Item Routes
	r.HandlerFunc(http.MethodGet, "/list/:lid/item", a.getItems)
	r.HandlerFunc(http.MethodPost, "/list/:lid/item", a.createItem)
	r.HandlerFunc(http.MethodGet, "/list/:lid/item/:iid", a.getItem)
	r.HandlerFunc(http.MethodPut, "/list/:lid/item/:iid", a.updateItem)
	r.HandlerFunc(http.MethodDelete, "/list/:lid/item/:iid", a.deleteItem)

	// Wrap the embedded handler in global middleware for logging
	a.Handler = web.RequestMW(r)

	return &a
}
