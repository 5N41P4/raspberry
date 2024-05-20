package main

import (
	// "errors"

	"net/http"
	"strconv"

	"github.com/5N41P4/raspberry/internal/data"
	"github.com/julienschmidt/httprouter"
)

func (app *application) getSchedules(w http.ResponseWriter, r *http.Request) {
	jobs := make([]data.Job, len(app.scheduler.Jobs))
	for i, job := range app.scheduler.Jobs {
		jobs[i] = *job
	}

	app.infoLog.Println("[SCHED]: GET Jobs")

	w.Header().Set("Content-Type", "application/json")
	err := app.writeJSON(w, http.StatusOK, jobs, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) addSchedule(w http.ResponseWriter, r *http.Request) {
	var input data.Job

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	app.infoLog.Printf("[SCHED]: Jobs\n")

	app.addJob(input)
}

func (app *application) deleteSchedule(w http.ResponseWriter, r *http.Request) {
	var cntx = httprouter.ParamsFromContext(r.Context()).ByName("id")

	id, err := strconv.Atoi(cntx)
	if err != nil {
		app.badRequestResponse(w, err)
	}

	app.infoLog.Printf("[SCHED]: Delete ID = %d", id)

	app.deleteJob(id)
}
