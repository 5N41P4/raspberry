package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/5N41P4/raspberry/internal/data"
)

// Interface Handlers

// Send the available interfaces as JSON
// getInterfaces handles the HTTP GET request to retrieve a list of interfaces.
// It populates the interfaces slice with data from the application's interfaces field,
// and writes the JSON representation of the interfaces to the response writer.
// If there are no interfaces available, the function returns early.
// If there is an error while writing the JSON response, the serverError method is called.
func (app *application) getInterfaces(w http.ResponseWriter, r *http.Request) {
	var interfaces []data.ApiInterface

	if app.interfaces == nil {
		app.writeJSON(w, http.StatusNoContent, nil, nil)
		app.errorLog.Println("no interfaces with monitoring capabilities given")
		return
	}

	for _, iface := range app.interfaces {
		apiiface := data.ApiInterface{
			Name:   iface.Name,
			State:  iface.State,
			Deauth: iface.Deauth.Running,
		}
		app.infoLog.Println(iface.State)
		interfaces = append(interfaces, apiiface)
	}

	if len(interfaces) == 0 || interfaces == nil {
		app.writeJSON(w, http.StatusNoContent, nil, nil)
		app.errorLog.Println("no interfaces with monitoring capabilities given")
		return
	}

	app.infoLog.Println("[INF] GET Interfaces")

	w.Header().Set("Content-Type", "application/json")
	err := app.writeJSON(w, http.StatusOK, interfaces, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

// Reveive JSON that requests an interface Action and tries to execute it.
// interfaceAction handles the API action for the interface.
// It reads the JSON input from the request, retrieves the corresponding interface,
// and tries to perform the requested action on the interface.
// If the action is successful, it logs the state of the interface.
// If there is an error, it logs the error message.
func (app *application) interfaceAction(w http.ResponseWriter, r *http.Request) {
	var input data.ApiAction

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	app.infoLog.Printf("%v", input)

	inf, ok := app.interfaces[input.Identifier]

	if !ok {
		app.errorLog.Println("[INF] Interface could not be found")
		app.badRequestResponse(w, errors.New("interface not found"))
		return
	}

	if input.Action == "stop" {
		inf.Stop()
		return
	}

	if inf.State != "up" {
		app.errorLog.Println("[INF] Requested bad interface action")
		app.badRequestResponse(w, errors.New("bad action requested"))
		return
	}

	inf.Target = getTarget(input.Target, &app.access, &app.clients)

	switch input.Action {
	case "capture":
		go inf.Capture(inf.Target)

	case "recon":
		go inf.Recon()

	default:
		app.errorLog.Println("[INF] Invalid action sent to interface")
		app.badRequestResponse(w, errors.New("invalid action"))
	}

	if input.Deauth {
		inf.Deauth.Running = true
		inf.Deauth.DeauthCan = make(chan struct{})
		go inf.RunDeauth(&app.access, &app.clients, inf.Target)
	}

	go inf.StopAfter(input.Time)
}

func getTarget(target string, access *map[string]*data.Accesspoint, clients *map[string]*data.Client) *data.Target {
	// If target is a client station then fill in the target information from the client
	cl, ok := (*clients)[target]
	if ok {
		return &data.Target{
			Bssid:   cl.Bssid,
			Station: cl.Station,
			Channel: strconv.Itoa((*access)[cl.Bssid].Channel),
		}
	}

	// If the target is a BSSID then fill in the target with the information from the accesspoint
	ap, ok := (*access)[target]
	if ok {
		return &data.Target{
			Bssid:   ap.Bssid,
			Station: "",
			Channel: strconv.Itoa(ap.Channel),
		}
	}

	// If the target could not be found then we fill in empty strings as a default
	return &data.Target{}
}
