package main

import (
	"net/http"
	"os"

	"github.com/5N41P4/raspberry/cmd/modules"
	"github.com/5N41P4/raspberry/internal/data"
)

// getCaptures handles the HTTP GET request to retrieve the list of capture files.
func (app *application) getFakeAP(w http.ResponseWriter, r *http.Request) {
	var output []data.ApiFakeAPStats

	files, err := os.ReadDir(modules.AttackBasePath)
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, file := range files {
		output = append(output, data.ApiFakeAPStats{
			Name:      file.Name(),
			Running:   false,
			Handshake: false,
			Key:       false,
		})
	}

	app.infoLog.Println("[Captures]")

	w.Header().Set("Content-Type", "application/json")
	err = app.writeJSON(w, http.StatusOK, output, nil)
	if err != nil {
		app.serverError(w, err)
	}
}
