package main

import (
	"net/http"

	"github.com/5N41P4/raspberry/cmd/modules"
	"github.com/5N41P4/raspberry/internal/data"
)

// System Information Handlers

func (app *application) systemInfo(w http.ResponseWriter, r *http.Request) {
	var system data.System
	var err error
	system.Disk, err = modules.GetDiskSpace()
	if err != nil {
		app.serverError(w, err)
		return
	}

	system.Cpu, err = modules.GetCpuUsage()
	if err != nil {
		app.serverError(w, err)
		return
	}

	system.Mem, err = modules.GetMemUsage()
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = app.writeJSON(w, http.StatusOK, system, nil)
	if err != nil {
		app.serverError(w, err)
	}

}
