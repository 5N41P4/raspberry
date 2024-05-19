package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/5N41P4/raspberry/internal/data"
)

// Sends JSON containing all the visible access points in the area
// getAP handles the HTTP GET request for retrieving access points.
// It retrieves the access points from the application's access slice,
// converts them to a slice of data.Accesspoint, and writes the JSON response
// containing the access points to the http.ResponseWriter.
func (app *application) getAP(w http.ResponseWriter, r *http.Request) {
	var aps []data.Accesspoint

	for _, ap := range app.access {
		aps = append(aps, *ap)
	}

	app.infoLog.Println("[APs]")

	w.Header().Set("Content-Type", "application/json")
	err := app.writeJSON(w, http.StatusOK, aps, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

// apAction handles the API actions related to the access point.
// It reads the JSON input from the request, performs the specified action,
// and sends the response back to the client.
// The available actions are: "reset", "delete", and "refresh".
// For "reset" action, it clears all the access points.
// For "delete" action, it removes the specified access point.
// For "refresh" action, it refreshes the access point lists.
// If the action is not found or there is an error reading the JSON input,
// it returns a bad request response.
func (app *application) apAction(w http.ResponseWriter, r *http.Request) {
	var input data.ApiAction

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
	app.infoLog.Printf("%s", input.Action)

	switch input.Action {
	case "reset":
		for k := range app.access {
			delete(app.access, k)
		}

	case "delete":
		delete(app.access, input.Identifier)

	case "refresh":
		app.refreshLists()

	default:
		app.badRequestResponse(w, errors.New("action not found"))
	}
}

// getSec retrieves the security information of the access points and sends it as a JSON response.
// It counts the number of access points using different security types (WEP, WPA, WPA2, WPA3)
// and returns the counts in a structured JSON format.
// The function takes a http.ResponseWriter and a http.Request as parameters.
// It writes the JSON response to the http.ResponseWriter and handles any errors that occur.
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
		WEP:  wep,
		WPA:  wpa,
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
