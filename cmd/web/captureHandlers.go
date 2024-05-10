package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/5N41P4/raspberry/internal/data"
)

// getCaptures handles the HTTP GET request to retrieve the list of capture files.
func (app *application) getCaptures(w http.ResponseWriter, r *http.Request) {
	var output data.ApiCaptures

	files, err := os.ReadDir("/usr/local/raspberry/captures")
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, file := range files {
		output.Files = append(output.Files, file.Name())
	}

	app.infoLog.Println("[Captures]")

	w.Header().Set("Content-Type", "application/json")
	err = app.writeJSON(w, http.StatusOK, output, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

// captureAction handles the capture action request.
// It reads the JSON input from the request body, validates it, and performs the corresponding action.
// If the action is "delete", it calls the captureDelete function with the provided identifier.
// Otherwise, it returns a bad request response with an error message.
func (app *application) captureAction(w http.ResponseWriter, r *http.Request) {
	var input data.ApiAction

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	app.infoLog.Printf("%s", input.Action)

	switch input.Action {

	case "delete":
		captureDelete(input.Identifier)

	default:
		app.badRequestResponse(w, errors.New("action not found"))
	}
}

// Helper functions

// deleta a Capture
func captureDelete(capName string) {
	os.RemoveAll("/usr/local/raspberry/captures/" + capName)
}
