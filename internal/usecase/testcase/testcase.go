package testcase

import (
	"context"
	"go-template-fiber/internal/svclogger"
)

func RunTask(aCtx context.Context, alog *svclogger.Log) {
	alog.Logger.Info().Msgf("Start task 'Test usecase'")

	_, cancel := context.WithCancel(aCtx)
	defer cancel()

	//
	alog.Logger.Info().Msgf("End task 'Test usecase'")
}
