package pgopher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
)

type Server struct {
	cfg Config
}

func NewServer(cfg Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Run(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go s.startScheduler(ctx, wg)

	mux := http.NewServeMux()

	mux.HandleFunc("/_ready", readinessProbe)
	mux.HandleFunc("/_live", livenessProbe)
	mux.HandleFunc("/api/v1/profile/", s.handleProfile)

	httpServer := &http.Server{
		Addr:    s.cfg.ListenAddress,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			slog.Error("failed to shutdown http server", slog.String("err", err.Error()))
		}
	}()

	slog.Info("starting http server", slog.String("listenAddr", s.cfg.ListenAddress))

	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server failed: %w", err)
	}

	wg.Wait()

	slog.Info("graceful shutdown completed, see you next time!")

	return nil
}
