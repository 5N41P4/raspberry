package main

import (
	"errors"
	"net/http"

	"github.com/5N41P4/raspberry/internal/data"
)

// Sends JSON containing all the visible clients
func (app *application) getClients(w http.ResponseWriter, r *http.Request) {
	var apicls []data.Client

	for _, cl := range app.clients {
		apicls = append(apicls, *cl)
	}

	app.infoLog.Println("[Clients]")

	w.Header().Set("Content-Type", "application/json")
	err := app.writeJSON(w, http.StatusOK, apicls, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) clientAction(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value("input").(*data.ApiSimpleAction)

	app.infoLog.Printf("%s", input.Action)

	switch input.Action {
	case "reset":
		for k := range app.clients {
			delete(app.clients, k)
		}

	case "delete":
		delete(app.clients, input.Identifier)

	case "refresh":
		app.refreshLists()

	default:
		app.badRequestResponse(w, errors.New("action not found"))
	}
}
