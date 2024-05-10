package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/5N41P4/raspberry/internal/data"
)

func (app *application) initScheduler() error {
	var jobPath string = "/usr/local/raspberry/jobs.json"
	var scheduler data.Scheduler

	// Read jobs from file
	file, err := os.Open(jobPath)
	if err != nil {
		app.errorLog.Println("error opening job file")
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(scheduler.Jobs)
	if err != nil {
		return fmt.Errorf("Error decoding JSON: %w", err)
	}

	app.scheduler = &scheduler

	return nil
}

func (app *application) addJob(job data.Job) error {
	return nil
}

func (app *application) deleteJob(id string) error {
	return nil
}
