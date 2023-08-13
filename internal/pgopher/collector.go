package pgopher

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type profileCollector struct {
	target ProfilingTarget
	logger slog.Logger
}

func (p profileCollector) Run() {
	p.logger.Info("collecting profile")

	resp, err := http.Get(fmt.Sprintf("%s?seconds=%d", p.target.URL, int(p.target.Duration.Seconds())))
	if err != nil {
		p.logger.Error("failed to collect profile", slog.String("error", err.Error()))
		return
	}

	defer resp.Body.Close()

	buf := &bytes.Buffer{}

	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		p.logger.Error("failed to read body", slog.String("error", err.Error()))
		return
	}
}
