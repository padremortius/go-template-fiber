package testcase

import (
	"context"

	"github.com/padremortius/go-template-fiber/pkgs/svclogger"
)

func RunTask(aCtx context.Context, alog *svclogger.Log) {
	alog.Infof("Start task 'Test usecase'")

	_, cancel := context.WithCancel(aCtx)
	defer cancel()

	//
	alog.Infof("End task 'Test usecase'")
}
