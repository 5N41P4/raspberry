package main

import (
	"net/http"
	"os"

	"github.com/5N41P4/raspberry/internal/data"
)

// Basic Server Functions

// home handles the HTTP request for the home page.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Println(wd)

	http.ServeFile(w, r, "/usr/local/raspberry/ui/dist/index.html")
}

// Test Handlers

func (app *application) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) apiTest(w http.ResponseWriter, r *http.Request) {
	data := data.Data{
		Labels: []string{"Data 1", "Data 2", "Data 3"},
		Datasets: []data.Dataset{
			{
				Label:           "Placeholder Data",
				Data:            []int{10, 20, 30},
				BackgroundColor: []string{"rgba(255, 99, 132, 0.3)", "rgba(54, 162, 235, 0.3)", "rgba(255, 206, 86, 0.3)"},
				BorderColor:     "white",
				BorderWidth:     1,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverError(w, err)
	}
}
