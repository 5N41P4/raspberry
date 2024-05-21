package main

import (
	"encoding/json"
	"os"

	"github.com/5N41P4/raspberry/cmd/modules"
	"github.com/5N41P4/raspberry/internal/data"
	"github.com/robfig/cron/v3"
)

// getScheduler returns a pointer to a data.Scheduler object populated with jobs from a JSON file.
func getScheduler() *data.Scheduler {
	var jobPath string = "/usr/local/raspberry/jobs.json"
	var scheduler data.Scheduler

	// Check if file exists
	if _, err := os.Stat(jobPath); os.IsNotExist(err) {
		// File does not exist, create it
		file, err := os.Create(jobPath)
		if err != nil {
			scheduler.Jobs = []*data.Job{}
			return &scheduler
		}
		defer file.Close()
	}

	// Read jobs from file
	file, err := os.Open(jobPath)
	if err != nil {
		scheduler.Jobs = []*data.Job{}
		return &scheduler
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	_ = decoder.Decode(&scheduler.Jobs)

	return &scheduler
}

// runScheduler runs the scheduler for the application.
// It iterates over the jobs in the scheduler and adds them to the cron scheduler.
// For each job, it checks if the corresponding interface exists in the application's interfaces map.
// If the interface is found, it adds a function to the cron scheduler that captures the target specified in the job.
// The function runs at the specified cron schedule.
func (app *application) runScheduler() {
	if app.scheduler.Proc != nil {
		app.scheduler.Proc.Stop()
	}

	app.scheduler.Proc = cron.New(cron.WithLogger(cron.PrintfLogger(app.infoLog)))
	app.scheduler.Proc.Start()
	if app.scheduler.Jobs == nil {
		return
	}

	for _, job := range app.scheduler.Jobs {
		inf, ok := app.interfaces[job.Cmd.Interface]
		if !ok {
			app.errorLog.Println("error: interface from job not found")
			continue
		}

		switch job.Cmd.Action {
		case "capture":
			fn := func(inf *modules.Interface, t *data.Target, delay int) func() {
				return func() {
					go inf.Capture(t)
					go inf.RunDeauth(&app.access, &app.clients, t)
					go inf.StopAfter(delay)
				}
			}(inf, job.Cmd.Target, job.Cmd.Time)
			_, err := app.scheduler.Proc.AddFunc(job.Cron, fn) // Call app.capture() separately
			if err != nil {
				app.errorLog.Println(err)
			}

		}
	}
}

func (app *application) stopScheduler() {
	app.scheduler.Proc.Stop()
	var jobPath string = "/usr/local/raspberry/jobs.json"

	file, err := os.Create(jobPath)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(app.scheduler.Jobs)
	if err != nil {
		app.errorLog.Println(err)
	}
}

// addJob adds a new job to the application's scheduler.
// It assigns a unique ID to the job and appends it to the list of jobs in the scheduler.
// After adding the job, it triggers the scheduler to run.
// Returns an error if there was a problem adding the job.
func (app *application) addJob(job data.Job) {
	job.ID = len(app.scheduler.Jobs)
	app.scheduler.Jobs = append(app.scheduler.Jobs, &job)
	app.runScheduler()
}

// deleteJob deletes a job from the application's scheduler based on the provided ID.
// If the ID matches a job in the scheduler, the job is removed from the scheduler.
// If the ID does not match any job, no action is taken.
func (app *application) deleteJob(id int) {
	if id < 0 || id >= len(app.scheduler.Jobs) {
		app.errorLog.Println("invalid job ID:", id)
		return
	}
	// Delete the job from the slice
	app.scheduler.Jobs = append(app.scheduler.Jobs[:id], app.scheduler.Jobs[id+1:]...)
	for id, job := range app.scheduler.Jobs {
		job.ID = id
	}
	app.runScheduler()
}
