package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mycreepy/pgopher/internal/pgopher"
	"gopkg.in/yaml.v3"
)

var (
	logLevel      = flag.String("log-level", "INFO", "log level")
	logJsonFormat = flag.Bool("log-json", false, "log in json format")
	cfgFile       = flag.String("config", "pgopher.yml", "config file")
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

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))

	if *logJsonFormat {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
	}

	slog.SetDefault(logger)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	file, err := os.Open(*cfgFile)
	if err != nil {
		slog.Error("failed to open config file", slog.String("err", err.Error()))
		os.Exit(1)
	}

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	cfg := pgopher.DefaultConfig()

	err = decoder.Decode(&cfg)
	if err != nil {
		slog.Error("failed to decode config", slog.String("err", err.Error()))
		os.Exit(1)
	}

	slog.Debug("loaded configuration", slog.String("file", *cfgFile), slog.Any("config", cfg))

	err = pgopher.NewServer(cfg).Run(ctx)
	if err != nil {
		slog.Error("failed to run server", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
