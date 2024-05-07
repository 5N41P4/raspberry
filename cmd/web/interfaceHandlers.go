package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/5N41P4/raspberry/internal/data"
)

// Interface Handlers

// Send the available interfaces as JSON
func (app *application) getInterfaces(w http.ResponseWriter, r *http.Request) {
	var interfaces []data.ApiInterface
	for _, iface := range app.interfaces {
		apiiface := data.ApiInterface{
			Name:  iface.Name,
			State: iface.State,
			Deauth: iface.Deauth,
		}
		app.infoLog.Println(iface.State)
		interfaces = append(interfaces, apiiface)
	}

	if len(interfaces) == 0 {
		return
	}

	app.infoLog.Println(interfaces)

	w.Header().Set("Content-Type", "application/json")
	err := app.writeJSON(w, http.StatusOK, interfaces, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

// Reveive JSON that requests an interface Action and tries to execute it.
func (app *application) interfaceAction(w http.ResponseWriter, r *http.Request) {
	var input data.ApiAction

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	app.infoLog.Printf("%s", input.Action)

	inf, ok := app.interfaces[input.Identifier]

	if !ok {
		app.errorLog.Println("Interface could not be found")
		return
	}

	state, err := inf.TryAction(input)

	if err != nil {
		app.errorLog.Println(err)
	}

	app.infoLog.Println(state)

}

// Handlers for interfacing with the Access Points and the Clients

// Sends JSON containing all the visible access points in the area
func (app *application) getAP(w http.ResponseWriter, r *http.Request) {
	var aps []data.ApiAP

	for _, ap := range app.access {
		apiap := data.ApiAP{
			Essid: ap.Essid,
			Bssid: ap.Bssid,
			Priv:  ap.Privacy,
		}
		aps = append(aps, apiap)
	}

	app.infoLog.Println("[APs]")

	w.Header().Set("Content-Type", "application/json")
	err := app.writeJSON(w, http.StatusOK, aps, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

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

// Sends JSON containing all the visible clients
func (app *application) getClients(w http.ResponseWriter, r *http.Request) {
	var apicls []data.ApiClient

	for _, cl := range app.clients {
		var bssid string
		if cl.Bssid == "(not associated) " {
			continue
		}

		// Translate the BSSID from the connected AP's into the more readable ESSID
		for apBssid, ap := range app.access {
			if apBssid == cl.Bssid && ap.Essid != "" {
				bssid = ap.Essid
				break
			}
			bssid = cl.Bssid
		}

		// Fill the Response with the appropriate strings
		apicl := data.ApiClient{
			Bssid:   bssid,
			Station: cl.Station,
		}
		apicls = append(apicls, apicl)
	}

	app.infoLog.Println("[Clients]")

	w.Header().Set("Content-Type", "application/json")
	err := app.writeJSON(w, http.StatusOK, apicls, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) clientAction(w http.ResponseWriter, r *http.Request) {
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
