package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/5N41P4/raspberry/cmd/modules"
	"github.com/5N41P4/raspberry/internal/data"
	"github.com/julienschmidt/httprouter"
)

// getCaptures handles the HTTP GET request to retrieve the list of capture files.
func (app *application) getCaptures(w http.ResponseWriter, r *http.Request) {
	var output data.ApiCaptures

	files, err := os.ReadDir(modules.CaptureBasePath)
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
	input := r.Context().Value("input").(*data.ApiAction)

	app.infoLog.Printf("%s", input.Action)

	switch input.Action {

	case "delete":
		captureDelete(input.Target)

	default:
		app.badRequestResponse(w, errors.New("action not found"))
	}
}

func (app *application) getCaptureWithId(w http.ResponseWriter, r *http.Request) {
	var id = httprouter.ParamsFromContext(r.Context()).ByName("id")
	path := modules.CaptureBasePath + "/" + id + "/-01.csv"

	aps, cls, err := modules.ParseCSV(path)
	if err != nil {
		app.badRequestResponse(w, err)
	}

	capture := data.ApiCapture{
		Accesspoints: aps,
		Clients:      cls,
	}

	app.infoLog.Printf("[CAPTURE] %s\n", id)

	w.Header().Set("Content-Type", "application/json")
	err = app.writeJSON(w, http.StatusOK, capture, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

// Helper functions

// deleta a Capture
func captureDelete(capName string) {
	os.RemoveAll("/usr/local/raspberry/captures/" + capName)
}
