package pgopher

import (
	"log/slog"
	"sync"

	"github.com/robfig/cron/v3"
)

func (s *Server) startScheduler(wg *sync.WaitGroup) {
	defer wg.Done()

	s.Logger.Info("starting scheduler", slog.Int("profilingTargets", len(s.cfg.ProfilingTargets)), slog.String("sink", s.cfg.Sink.Type))

	scheduler := cron.New()

	for _, target := range s.cfg.ProfilingTargets {
		logger := s.Logger.With(slog.String("target", target.Name))

		_, err := scheduler.AddJob(target.Schedule, profileCollector{
			ctx:        s.Context,
			logger:     *logger,
			target:     target,
			sink:       s.cfg.Sink,
			s3Client:   s.s3Client,
			kubeClient: s.kubeClient,
		})
		if err != nil {
			logger.Error("failed to create collector for profiling target", slog.String("err", err.Error()))
			continue
		}
	}

	go func() {
		<-s.Context.Done()
		<-scheduler.Stop().Done()
	}()

	scheduler.Run()
}
