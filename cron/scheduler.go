package cron

import (
	"context"

	"github.com/ItsNotGoodName/smtpbridge/pkg/sutureext"
	"github.com/reugn/go-quartz/quartz"
)

type Scheduler struct {
	*quartz.StdScheduler
}

func NewScheduler() Scheduler {
	return Scheduler{
		StdScheduler: quartz.NewStdSchedulerWithOptions(quartz.StdSchedulerOptions{}),
	}
}

func (s Scheduler) Serve(ctx context.Context) error {
	s.StdScheduler.Start(ctx)

	<-ctx.Done()

	s.StdScheduler.Wait(context.Background())

	return nil
}

// ScheduleJob is used to defer job scheduling after the Scheduler has been started by Suture.
// Any error will kill the parent supervisor tree.
func ScheduleJob(fn func(context.Context) error) sutureext.ServiceFunc {
	return sutureext.NewServiceFunc("cron.ScheduleJob", sutureext.OneShotFunc(fn))
}
