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
		Action: data.ApiAction{
			Identifier: "wlan1",
			Action:     "capture",
			Target:     "1C:91:80:D1:D8:BE",
			Deauth:     false,
		},
		Cron: "10 * * * 5",
	}

	job1 := data.Job{
		Action: data.ApiAction{
			Identifier: "wlan1",
			Action:     "capture",
			Target:     "6C:FF:CE:EC:15:94",
			Deauth:     false,
		},
		Cron: "0 0 * * *",
	}

	app := &application{
		interfaces: map[string]*modules.Interface{
			"wlan0": &wlan0,
			"wlan1": &wlan1,
		},
		access:    make(map[string]*data.AppAP),
		clients:   make(map[string]*data.AppClient),
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

	// Check job0
	if app.scheduler.Jobs[0].Action.Identifier != job0.Action.Identifier {
		t.Errorf("Expected job0 identifier %s, got %s", job0.Action.Identifier, app.scheduler.Jobs[0].Action.Identifier)
	}
	if app.scheduler.Jobs[0].Action.Action != job0.Action.Action {
		t.Errorf("Expected job0 action %s, got %s", job0.Action.Action, app.scheduler.Jobs[0].Action.Action)
	}
	if app.scheduler.Jobs[0].Action.Target != job0.Action.Target {
		t.Errorf("Expected job0 target %s, got %s", job0.Action.Target, app.scheduler.Jobs[0].Action.Target)
	}
	if app.scheduler.Jobs[0].Action.Deauth != job0.Action.Deauth {
		t.Errorf("Expected job0 deauth %t, got %t", job0.Action.Deauth, app.scheduler.Jobs[0].Action.Deauth)
	}
	if app.scheduler.Jobs[0].Cron != job0.Cron {
		t.Errorf("Expected job0 cron %s, got %s", job0.Cron, app.scheduler.Jobs[0].Cron)
	}

	// Check job1
	if app.scheduler.Jobs[1].Action.Identifier != job1.Action.Identifier {
		t.Errorf("Expected job1 identifier %s, got %s", job1.Action.Identifier, app.scheduler.Jobs[1].Action.Identifier)
	}
	if app.scheduler.Jobs[1].Action.Action != job1.Action.Action {
		t.Errorf("Expected job1 action %s, got %s", job1.Action.Action, app.scheduler.Jobs[1].Action.Action)
	}
	if app.scheduler.Jobs[1].Action.Target != job1.Action.Target {
		t.Errorf("Expected job1 target %s, got %s", job1.Action.Target, app.scheduler.Jobs[1].Action.Target)
	}
	if app.scheduler.Jobs[1].Action.Deauth != job1.Action.Deauth {
		t.Errorf("Expected job1 deauth %t, got %t", job1.Action.Deauth, app.scheduler.Jobs[1].Action.Deauth)
	}
	if app.scheduler.Jobs[1].Cron != job1.Cron {
		t.Errorf("Expected job1 cron %s, got %s", job1.Cron, app.scheduler.Jobs[1].Cron)
	}

	// Check function in app.scheduler.Proc.Jobs
	if len(app.scheduler.Proc.Entries()) != len(app.scheduler.Jobs) {
		t.Errorf("Expected app.scheduler.Proc.Entries to be %d, got %d", len(app.scheduler.Jobs), len(app.scheduler.Proc.Entries()))
	}
}
