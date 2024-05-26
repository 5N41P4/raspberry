package main

import (
	"errors"
	"net/http"

	"github.com/5N41P4/raspberry/cmd/modules"
	"github.com/5N41P4/raspberry/internal/data"
)

const (
	Up modules.InterfaceState = iota
	Inet
	Recon
	Capture
	AccessPoint
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

func (app *application) interfaceAction(w http.ResponseWriter, r *http.Request) {
	input := r.Context().Value("input").(*data.ApiInterfaceAction)
	inf := r.Context().Value("inf").(*modules.Interface)

	if input.Action == "stop" {
		inf.Stop()
		return
	}

	if inf.State != modules.InterfaceStates[Up] {
		app.errorLog.Println("[INF] Requested bad interface action")
		app.badRequestResponse(w, errors.New("bad action requested"))
		return
	}

	inf.Target = &input.Target

	switch input.Action {
	case modules.InterfaceStates[AccessPoint]:
		ap, err := modules.NewFakeAP(inf.Name, inf.Target)
		if err != nil {
			app.errorLog.Println("[INF] Could not create fake AP")
			app.serverError(w, err)
			return
		}
		inf.FakeAP = ap
		inf.State = modules.InterfaceStates[AccessPoint]
		go inf.FakeAP.Start()

	case modules.InterfaceStates[Capture]:
		go inf.Capture(inf.Target)

	case modules.InterfaceStates[Recon]:
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
