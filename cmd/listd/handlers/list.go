package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *Application) getLists(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (a *Application) createList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (a *Application) getList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	listID := ps.ByName("lid")
}

func (a *Application) updateList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	listID := ps.ByName("lid")
}

func (a *Application) deleteList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	listID := ps.ByName("lid")
}
