package main

import (
	"fmt"
	"net/http"

	"github.com/5N41P4/raspberry/internal/data"
)

// Handlers for the filter lists:

// Get handler for the filter lists:
func (app *application) getFilters(w http.ResponseWriter, r *http.Request) {
	list := r.Context().Value("list").(*Filter)
	w.Header().Set("Content-Type", "application/json")

	err := app.writeJSON(w, http.StatusOK, list, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

// Post handler
func (app *application) filterAction(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value("input").(*data.ApiSimpleAction)
	list := r.Context().Value("list").(*Filter)

	fmt.Fprintf(w, "%+v\n", input)

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
