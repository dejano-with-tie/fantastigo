package server

import (
	"github.com/dejano-with-tie/fantastigo/internal/common/apperr"
	"github.com/dejano-with-tie/fantastigo/internal/ector/app"
	"github.com/dejano-with-tie/fantastigo/internal/ector/metrics"
	"github.com/labstack/echo/v4"
	"net/http"
)

type EctorHttpHandler struct {
	app *app.Ector
}

var _ ServerInterface = &EctorHttpHandler{}

func NewEctorHttpHandler(app *app.Ector) *EctorHttpHandler {
	return &EctorHttpHandler{app: app}
}

func (e EctorHttpHandler) GetIdentity(ctx echo.Context) error {
	i, err := e.app.GetIdentity()
	if err != nil {
		return apperr.Wrap(apperr.ErrCodeInternal, err)
	}

	return ctx.JSON(http.StatusOK, i)
}

func (e EctorHttpHandler) GetStatus(ctx echo.Context) error {
	info, err := e.app.Status()
	if err != nil {
		return err
	}

	r := VehicleStatusResponse{Vin: info.Vin}

	return ctx.JSON(http.StatusOK, &r)
}

func (e EctorHttpHandler) GetMetrics(ctx echo.Context, params GetMetricsParams) error {
	// wildcard matches all metrics
	queryName := "*"
	if params.MeasurementName != nil {
		queryName = *params.MeasurementName
	}

	q := metrics.Query{
		MeasurementName: queryName,
	}

	m, err := e.app.GetMetrics(q)
	if err != nil {
		return apperr.Wrap(apperr.ErrCodeInternal, err)
	}

	return ctx.JSON(http.StatusOK, m)
}
