package data

import "github.com/robfig/cron/v3"

type Scheduler struct {
	Jobs []*Job `json:"jobs"`
	Proc *cron.Cron
}

type Job struct {
	ID     string    `json:"id"`
	Cron   string    `json:"cron"`
	Action ApiAction `json:"action"`
}
