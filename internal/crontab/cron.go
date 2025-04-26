package crontab

import (
	"context"
	"go-template-fiber/internal/usecase/testcase"
)

func (ct *Crontab) LoadTasks(aCtx context.Context, opts *CronOpts) {
	ctx, cancel := context.WithCancel(aCtx)
	defer cancel()

	taskCount := 0
	for _, job := range opts.Jobs {
		if !job.Disable {
			taskCount++
		}
	}
	if taskCount > 0 {
		ct.Logger.Logger.Debug().Msgf("taskCount = %v", taskCount)
		ct.WGroup.Add(taskCount)
		if !opts.Jobs[0].Disable {
			ct.Logger.Logger.Info().Msgf("Add new task. { Name: %v, Schedule: %v }", opts.Jobs[0].Name, opts.Jobs[0].Schedule)
			_, _ = ct.CronServer.AddFunc(opts.Jobs[0].Schedule, func() {
				testcase.RunTask(ctx, ct.Logger)
			})
		}
	}
}
