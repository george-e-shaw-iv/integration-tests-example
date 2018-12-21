package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *Application) getItems(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	listID := ps.ByName("lid")
}

func (a *Application) createItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	listID := ps.ByName("lid")
}

func (a *Application) getItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	listID := ps.ByName("lid")
	itemID := ps.ByName("iid")
}

func (a *Application) updateItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	listID := ps.ByName("lid")
	itemID := ps.ByName("iid")
}

func (a *Application) deleteItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	listID := ps.ByName("lid")
	itemID := ps.ByName("iid")
}
