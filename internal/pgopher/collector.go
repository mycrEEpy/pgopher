package pgopher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type profileCollector struct {
	ctx    context.Context
	logger slog.Logger
	target ProfilingTarget
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
}
