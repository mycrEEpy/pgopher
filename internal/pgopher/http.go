package pgopher

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
)

func readinessProbe(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func livenessProbe(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func (s *Server) handleProfile(ctx echo.Context) error {
	profile := ctx.Param("profile")

	if len(profile) == 0 || strings.Contains(profile, "..") {
		slog.Error("invalid profile", slog.String("profile", profile))
		return ctx.NoContent(http.StatusBadRequest)
	}

	switch s.cfg.Sink.Type {
	case "file":
		return ctx.File(filepath.Join(s.cfg.Sink.FileSinkOptions.Folder, profile))
	case "s3":
		resp, err := s.s3Client.GetObject(ctx.Request().Context(), &s3.GetObjectInput{
			Bucket: aws.String(s.cfg.Sink.S3SinkOptions.Bucket),
			Key:    aws.String(fmt.Sprintf("profile=%s/%s.pgo", profile, profile)),
		})
		if err != nil {
			slog.Error("failed to get profile from s3 sink", slog.String("err", err.Error()), slog.String("profile", profile))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		defer resp.Body.Close()

		file, err := os.CreateTemp(os.TempDir(), "pgopher-*.pgo")
		if err != nil {
			slog.Error("failed to create temporary file", slog.String("err", err.Error()), slog.String("profile", profile))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		defer file.Close()
		defer os.Remove(file.Name())

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			slog.Error("failed to write temporary file", slog.String("err", err.Error()), slog.String("profile", profile))
			return ctx.NoContent(http.StatusInternalServerError)
		}

		return ctx.File(file.Name())
	default:
		return ctx.NoContent(http.StatusInternalServerError)
	}
}
