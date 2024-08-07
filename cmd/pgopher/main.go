package main

import (
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/mycreepy/pgopher/internal/pgopher"
)

var (
	cfgFile      = flag.String("config", "pgopher.yml", "config file")
	pprofEnabled = flag.Bool("pprof", false, "enable pprof endpoint")
)

func main() {
	flag.Parse()

	cfg, err := pgopher.LoadConfig(*cfgFile)
	if err != nil {
		slog.Error("failed to load config", slog.String("err", err.Error()))
		os.Exit(1)
	}

	s, err := pgopher.NewServer(cfg)
	if err != nil {
		slog.Error("failed to create pgopher server", slog.String("err", err.Error()))
		os.Exit(1)
	}

	defer s.CancelContext()

	if *pprofEnabled {
		go pprofServer(cfg.PprofListenAddress, s.Logger)
	}

	err = s.Run()
	if err != nil {
		slog.Error("failed to run pgopher server", slog.String("err", err.Error()))
		os.Exit(1)
	}
}

func pprofServer(listenAddr string, logger *slog.Logger) {
	pprofMux := http.NewServeMux()
	pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)

	server := &http.Server{Addr: listenAddr, Handler: pprofMux}

	logger.Info("starting pprof server", slog.String("listenAddr", listenAddr))

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("pprof server failed", slog.String("err", err.Error()))
	}
}
