package pgopher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type profileCollector struct {
	ctx    context.Context
	logger slog.Logger
	target ProfilingTarget
	sink   Sink
}

func (p profileCollector) Run() {
	p.logger.Info("collecting profile")

	url := fmt.Sprintf("%s?seconds=%d", p.target.URL, int(p.target.Duration.Seconds()))

	req, err := http.NewRequestWithContext(p.ctx, http.MethodGet, url, nil)
	if err != nil {
		p.logger.Error("failed to create http request", slog.String("err", err.Error()))
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		p.logger.Error("failed to collect profile", slog.String("err", err.Error()))
		return
	}

	defer resp.Body.Close()

	buf := &bytes.Buffer{}

	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		p.logger.Error("failed to read body", slog.String("err", err.Error()))
		return
	}

	switch p.sink.Type {
	case "file":
		filePath := filepath.Join(p.sink.FileSinkOptions.Folder, fmt.Sprintf("%s.pgo", p.target.Name))

		file, err := os.Create(filePath)
		if err != nil {
			p.logger.Error("failed to create file for sink", slog.String("err", err.Error()))
			return
		}

		defer file.Close()

		_, err = file.Write(buf.Bytes())
		if err != nil {
			p.logger.Error("failed to write to sink file", slog.String("err", err.Error()), slog.String("file", file.Name()))
			return
		}
	default:
		p.logger.Error("unknown sink", slog.String("sink", p.sink.Type))
		return
	}
}
