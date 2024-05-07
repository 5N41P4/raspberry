package main

import (
	"fmt"
	"net/http"

	"github.com/5N41P4/raspberry/internal/data"
	"github.com/julienschmidt/httprouter"
)

// Handlers for the filter lists:

// Get handler for the filter lists:
func (app *application) getFilters(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")

	w.Header().Set("Content-Type", "application/json")

	if id == "clients" {
		app.infoLog.Println("[ClientFilter]")
		err := app.writeJSON(w, http.StatusOK, app.filters.ClFilter, nil)
		if err != nil {
			app.serverError(w, err)
		}
		return
	}
	if id == "aps" {
		app.infoLog.Println("[ApFilter]")
		err := app.writeJSON(w, http.StatusOK, app.filters.ApFilter, nil)
		if err != nil {
			app.serverError(w, err)
		}
		return
	}
}

// Post handler
func (app *application) filterAction(w http.ResponseWriter, r *http.Request) {
	var input data.ApiAction
	var list *Filter

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)

	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if id == "clients" {
		list = app.filters.ClFilter
	} else if id == "aps" {
		list = app.filters.ApFilter
	} else {
		app.clientError(w, http.StatusBadRequest)
	}

	switch input.Action {
	case "switch":
		list.Switch()

	case "add":
		list.Add(input.Identifier)

	case "delete":
		list.Delete(input.Identifier)

	case "reset":
		list.Reset()
	}

	app.refreshLists()
}
