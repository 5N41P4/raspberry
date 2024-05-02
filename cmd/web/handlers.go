package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/5N41P4/rpine/internal/data"
)

// Basic Server Functions

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Println(wd)

	http.ServeFile(w, r, "/usr/local/raspberry/ui/dist/index.html")
}

// Handlers for the security API chart

func (app *application) getSec(w http.ResponseWriter, r *http.Request) {
  var output data.ApiSecurity
  var wep, wpa, wpa2, wpa3 int

	for _, ap := range app.access {
    secTypes := strings.Fields(ap.Privacy)
    for _, sec := range secTypes {
      switch sec {
      case "WEP":
        wep += 1
      case "WPA":
        wpa += 1
      case "WPA2":
        wpa2 += 1
      case "WPA3": 
        wpa3 += 1
      } 
    }
	}

  output = data.ApiSecurity{
    WEP: wep,
    WPA: wpa,
    WPA2: wpa2,
    WPA3: wpa3,
  }

	app.infoLog.Println("[Security]")

	w.Header().Set("Content-Type", "application/json")
	err := app.writeJSON(w, http.StatusOK, output, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

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
