package pgopher

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func readinessProbe(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func livenessProbe(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}
