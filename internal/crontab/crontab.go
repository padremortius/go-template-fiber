package crontab

import (
	"context"
	"sync"

	"github.com/padremortius/go-template-fiber/internal/svclogger"

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
		Logger     *svclogger.Log
		WGroup     *sync.WaitGroup
	}
)

func (ct Crontab) StartCron() {
	ct.Logger.Logger.Info().Msg("Start crontab")
	ct.CronServer.Start()
	for i, item := range ct.CronServer.Entries() {
		ct.Logger.Logger.Info().Msgf("Task %v next time start %v", i, item.Next)
	}
}

func (ct *Crontab) StopCron() {
	ct.Logger.Logger.Info().Msg("Waiting for stop crontab")
	ct.CronServer.Stop()
	ct.WGroup.Done()
	ct.Logger.Logger.Info().Msg("Stop crontab")
}

func New(aCtx context.Context, alog *svclogger.Log, opts *CronOpts) Crontab {
	var ct *cron.Cron
	if opts.EnableSeconds {
		ct = cron.New(cron.WithSeconds())
	} else {
		ct = cron.New()
	}

	return Crontab{
		CronServer: ct,
		Ctx:        aCtx,
		Logger:     alog,
		WGroup:     &sync.WaitGroup{},
	}
}
