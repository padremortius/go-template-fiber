package crontab

import (
	"context"
	"sync"

	cron "github.com/robfig/cron/v3"
)

type (
	CronOpts struct {
		EnableSeconds bool  `yaml:"enableSeconds" json:"enableSeconds" validate:"required"`
		Jobs          []Job `yaml:"jobs" json:"jobs" validate:"required"`
	}

	Job struct {
		Name     string `yaml:"name" json:"name" validate:"required"`
		Schedule string `yaml:"schedule" json:"schedule" validate:"required"`
		Disable  bool   `yaml:"disable" json:"disable" validate:"required"`
	}

	Crontab struct {
		CronServer *cron.Cron
		Ctx        context.Context
		WGroup     *sync.WaitGroup
	}
)

func (ct Crontab) StartCron() {
	ct.CronServer.Start()
}

func (ct *Crontab) StopCron() {
	ct.CronServer.Stop()
	ct.WGroup.Done()
}

func New(aCtx context.Context, opts *CronOpts) Crontab {
	var ct *cron.Cron
	if opts.EnableSeconds {
		ct = cron.New(cron.WithSeconds())
	} else {
		ct = cron.New()
	}

	return Crontab{
		CronServer: ct,
		Ctx:        aCtx,
		WGroup:     &sync.WaitGroup{},
	}
}
