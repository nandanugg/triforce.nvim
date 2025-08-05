package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewEchoServer(openapiBytes []byte) (*echo.Echo, error) {
	e := echo.New()

	var err error
	e.Binder, err = newEchoBinder(openapiBytes)
	if err != nil {
		return nil, fmt.Errorf("laod openapi schema: %w", err)
	}

	e.HideBanner = true
	e.Use(
		middleware.Recover(),
		middleware.RequestID(),
		newLogMiddleware(),
		newPrometheusMiddleware(),
	)
	e.Add(http.MethodGet, "/health", func(echo.Context) error { return nil })
	e.Add(http.MethodGet, "/metrics", echoprometheus.NewHandler())

	return e, nil
}

// StartEchoServer start e in specified port.
func StartEchoServer(e *echo.Echo, port uint) error {
	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	startCh := make(chan error)
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
			startCh <- err
		}
	}()

	select {
	case <-sigCtx.Done():
		slog.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return e.Shutdown(ctx)
	case err := <-startCh:
		return err
	}
}
