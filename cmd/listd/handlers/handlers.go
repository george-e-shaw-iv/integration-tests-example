package handlers

import (
	"net/http"

	"github.com/george-e-shaw-iv/integration-tests-example/cmd/listd/configuration"
	"github.com/george-e-shaw-iv/integration-tests-example/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

// Application is the struct that contains the server handler as well as
// any references to services that the application needs.
type Application struct {
	db      *sqlx.DB
	cfg     *configuration.Config
	handler http.Handler
}

// NewApplication returns a new pointer to Application with route definitions
// initiated.
func NewApplication(db *sqlx.DB, cfg *configuration.Config) *Application {
	a := &Application{
		db:  db,
		cfg: cfg,
	}
	a.initHandlers()

	return a
}

// initHandlers initiates the routes attached to the server handler within the
// Application type.
func (a *Application) initHandlers() {
	r := httprouter.New()

	// Kubernetes Probes
	probeHandler := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
	}

	r.GET("/ready", probeHandler)
	r.GET("/healthy", probeHandler)

	// List Routes
	r.GET("/list", a.getLists)
	r.POST("/list", a.createList)
	r.GET("/list/:lid", a.getList)
	r.PUT("/list/:lid", a.updateList)
	r.DELETE("/list/:lid", a.deleteList)

	// Item Routes
	r.GET("/list/:lid/item", a.getItems)
	r.POST("/list/:lid/item", a.createItem)
	r.GET("/list/:lid/item/:iid", a.getItem)
	r.PUT("/list/:lid/item/:iid", a.updateItem)
	r.DELETE("/list/:lid/item/:iid", a.deleteItem)

	a.handler = web.RequestMW(r)
}

// ServeHTTP implements the http handler interface for type Application.
func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}
