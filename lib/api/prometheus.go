package api

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

func newPrometheusMiddleware() echo.MiddlewareFunc {
	return echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		Registerer: prometheus.NewRegistry(),
		Subsystem:  "http",
		LabelFuncs: map[string]echoprometheus.LabelValueFunc{
			"url": func(c echo.Context, _ error) string { // overrides default 'url' label value
				if c.Path() == "" {
					return "NOT_FOUND"
				}
				return c.Path()
			},
		},
	})
}
