package data

import (
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	Jobs []*Job `json:"jobs"`
	Proc *cron.Cron
}

type Job struct {
	ID   int    `json:"id"`
	Cron string `json:"cron"`
	Cmd  Cmd    `json:"cmd"`
}

type Cmd struct {
	Interface string  `json:"interface"`
	Time      int     `json:"time"`
	Action    string  `json:"action"`
	Target    *Target `json:"target"`
	Deauth    bool    `json:"deauth"`
}
