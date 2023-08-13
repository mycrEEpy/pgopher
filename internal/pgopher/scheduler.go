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
			target: target,
			logger: *logger,
		})
		if err != nil {
			logger.Error("failed to create collector for profiling target", slog.String("error", err.Error()))
			continue
		}
	}

	go func() {
		<-ctx.Done()
		<-scheduler.Stop().Done()
	}()

	scheduler.Run()
}
