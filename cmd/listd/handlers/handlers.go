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
	DB      *sqlx.DB
	handler http.Handler
}

// ServeHTTP implements the http.Handler interface for the Application type.
func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}

// NewApplication returns a new pointer to Application with route definitions
// initiated.
func NewApplication(db *sqlx.DB) *Application {
	a := Application{
		DB: db,
	}

	router := httprouter.New()

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
	router.HandlerFunc(http.MethodGet, "/ready", probeHandler)
	router.HandlerFunc(http.MethodGet, "/healthy", probeHandler)

	// List Routes
	router.HandlerFunc(http.MethodGet, "/list", a.getLists)
	router.HandlerFunc(http.MethodPost, "/list", a.createList)
	router.HandlerFunc(http.MethodGet, "/list/:lid", a.getList)
	router.HandlerFunc(http.MethodPut, "/list/:lid", a.updateList)
	router.HandlerFunc(http.MethodDelete, "/list/:lid", a.deleteList)

	// Item Routes
	router.HandlerFunc(http.MethodGet, "/list/:lid/item", a.getItems)
	router.HandlerFunc(http.MethodPost, "/list/:lid/item", a.createItem)
	router.HandlerFunc(http.MethodGet, "/list/:lid/item/:iid", a.getItem)
	router.HandlerFunc(http.MethodPut, "/list/:lid/item/:iid", a.updateItem)
	router.HandlerFunc(http.MethodDelete, "/list/:lid/item/:iid", a.deleteItem)

	// Wrap the router in middleware used for logging requests and set the application
	// handler to utilize the returned http.Handler from RequestMW.
	a.handler = web.RequestMW(router)

	return &a
}
