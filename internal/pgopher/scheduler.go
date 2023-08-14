package pgopher

import (
	"context"
	"log/slog"
	"sync"

	"github.com/robfig/cron/v3"
)

func (s *Server) startScheduler(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	slog.Info("starting scheduler", slog.Int("profilingTargets", len(s.cfg.ProfilingTargets)))

	scheduler := cron.New()

	for name, target := range s.cfg.ProfilingTargets {
		logger := slog.With(slog.String("target", name))

		_, err := scheduler.AddJob(target.Schedule, profileCollector{
			ctx:    ctx,
			logger: *logger,
			target: target,
		})
		if err != nil {
			logger.Error("failed to create collector for profiling target", slog.String("err", err.Error()))
			continue
		}
	}

	go func() {
		<-ctx.Done()
		<-scheduler.Stop().Done()
	}()

	scheduler.Run()
}
