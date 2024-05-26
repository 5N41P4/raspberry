package main

import (
	"log"
	"testing"

	"github.com/5N41P4/raspberry/cmd/modules"
	"github.com/5N41P4/raspberry/internal/data"
)

func TestRunScheduler(t *testing.T) {
	wlan0 := modules.Interface{
		Name:  "wlan0",
		State: "up",
	}

	wlan1 := modules.Interface{
		Name:  "wlan1",
		State: "recon",
	}

	job0 := data.Job{
		Cmd: data.Cmd{
			Interface: "wlan0",
			Action:    "capture",
			Time:      5,
			Target: &data.Target{
				Bssid:   "6C:FF:CE:EC:15:94",
				Essid:   "Test",
				Station: "FF:FF:FF:FF:FF:FF",
				Channel: "1",
				Privacy: "WPA2",
			},
			Deauth: false,
		},
		Cron: "10 * * * 5",
	}

	job1 := data.Job{
		Cmd: data.Cmd{
			Interface: "wlan0",
			Action:    "capture",
			Time:      5,
			Target: &data.Target{
				Bssid:   "5B:AB:CD:EF:12:34",
				Essid:   "Top_5G",
				Station: "",
				Channel: "36",
				Privacy: "WPA",
			},
			Deauth: false,
		},
		Cron: "10 * * * 5",
	}

	app := &application{
		interfaces: map[string]*modules.Interface{
			"wlan0": &wlan0,
			"wlan1": &wlan1,
		},
		access:    make(map[string]*data.Accesspoint),
		clients:   make(map[string]*data.Client),
		filters:   newFilterList(),
		scheduler: getScheduler(),
	}

	app.runScheduler()

	aps, cls, err := modules.ParseCSV("/Users/aurel/Source/raspberry/testfiles/try-02.csv")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	for _, ap := range aps {
		if app.filters.ApFilter.IsAllowed(ap.Bssid) {
			app.access[ap.Bssid] = &ap
		}
	}

	for _, cli := range cls {
		if app.filters.ClFilter.IsAllowed(cli.Station) {
			app.clients[cli.Station] = &cli
		}
	}

	app.addJob(job0)
	app.addJob(job1)

	// Assert that the jobs are added correctly
	if len(app.scheduler.Jobs) != 2 {
		t.Errorf("Expected 2 jobs, got %d", len(app.scheduler.Jobs))
	}

	for _, job := range app.scheduler.Jobs {
		if job.Cmd.Interface != "wlan0" {
			t.Errorf("Expected job interface to be wlan0, got %s", job.Cmd.Interface)
		}
		if job.Cron != "10 * * * 5" {
			t.Errorf("Expected job cron to be 10 * * * 5, got %s", job.Cron)
		}
	}

	// Check function in app.scheduler.Proc.Jobs
	if len(app.scheduler.Proc.Entries()) != len(app.scheduler.Jobs) {
		t.Errorf("Expected app.scheduler.Proc.Entries to be %d, got %d", len(app.scheduler.Jobs), len(app.scheduler.Proc.Entries()))
	}
}
