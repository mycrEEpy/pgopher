package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/mycreepy/pgopher/internal/pgopher"
)

var (
	logLevel      = flag.String("log-level", "INFO", "log level")
	logJsonFormat = flag.Bool("log-json", false, "log in json format")
	logSource     = flag.Bool("log-source", false, "log source code position")
	cfgFile       = flag.String("config", "pgopher.yml", "config file")
	pprofEnabled  = flag.Bool("pprof", false, "enable pprof endpoint")
)

func init() {
	flag.Parse()
	setupLogger()
}

func setupLogger() {
	level := &slog.LevelVar{}

	err := level.UnmarshalText([]byte(*logLevel))
	if err != nil {
		slog.Error("failed to parse log level", slog.String("err", err.Error()))
		os.Exit(1)
	}

	handler := &slog.HandlerOptions{
		AddSource: *logSource,
		Level:     level,
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, handler))

	if *logJsonFormat {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, handler))
	}

	slog.SetDefault(logger)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	cfg, err := pgopher.LoadConfig(*cfgFile)
	if err != nil {
		slog.Error("failed to load config", slog.String("err", err.Error()))
		os.Exit(1)
	}

	slog.Debug("loaded configuration", slog.String("file", *cfgFile), slog.Any("config", cfg))

	if *pprofEnabled {
		go pprofServer(cfg.PprofListenAddress)
	}

	s, err := pgopher.NewServer(cfg)
	if err != nil {
		slog.Error("failed to create pgopher server", slog.String("err", err.Error()))
		os.Exit(1)
	}

	err = s.Run(ctx)
	if err != nil {
		slog.Error("failed to run pgopher server", slog.String("err", err.Error()))
		os.Exit(1)
	}
}

func pprofServer(listenAddr string) {
	pprofMux := http.NewServeMux()
	pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)

	server := &http.Server{Addr: listenAddr, Handler: pprofMux}

	slog.Info("starting pprof server", slog.String("listenAddr", listenAddr))

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("pprof server failed", slog.String("err", err.Error()))
	}
}
