package cron

import (
	"context"

	"github.com/padremortius/go-template-fiber/internal/usecase/testcase"
	"github.com/padremortius/go-template-fiber/pkgs/crontab"
	"github.com/padremortius/go-template-fiber/pkgs/svclogger"
)

func LoadTasks(aCtx context.Context, ct crontab.Crontab, opts *crontab.CronOpts, log *svclogger.Log) {
	ctx, cancel := context.WithCancel(aCtx)
	defer cancel()

	taskCount := 0
	for _, job := range opts.Jobs {
		if !job.Disable {
			taskCount++
		}
	}
	if taskCount > 0 {
		log.Debugf("taskCount = %v", taskCount)
		if !opts.Jobs[0].Disable {
			log.Infof("Add new task. { Name: %v, Schedule: %v }", opts.Jobs[0].Name, opts.Jobs[0].Schedule)
			_, err := ct.CronServer.AddFunc(opts.Jobs[0].Schedule, func() {
				ct.WGroup.Add(1)
				defer ct.WGroup.Done()
				testcase.RunTask(ctx, log)
			})
			if err != nil {
				log.Errorf("Error set crontask. Error: %v", err)
				ct.WGroup.Done()
			}
		}
	}
}
