package main

import (
	"net/http"

	"github.com/5N41P4/raspberry/cmd/modules"
)

// System Information Handlers

func (app *application) diskUsage(w http.ResponseWriter, r *http.Request) {
	diskUsage, err := modules.GetDiskSpace()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Println("[Disk]")

	w.Header().Set("Content-Type", "application/json")
	err = app.writeJSON(w, http.StatusOK, diskUsage, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) cpuUsage(w http.ResponseWriter, r *http.Request) {
	cpuUsage, err := modules.GetCpuUsage()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Println("[CPU]")

	w.Header().Set("Content-Type", "application/json")
	err = app.writeJSON(w, http.StatusOK, cpuUsage, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) memUsage(w http.ResponseWriter, r *http.Request) {
	memUsage, err := modules.GetMemUsage()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Println("[Memory]")

	w.Header().Set("Content-Type", "appliaction/json")
	err = app.writeJSON(w, http.StatusOK, memUsage, nil)
	if err != nil {
		app.serverError(w, err)
	}
}
