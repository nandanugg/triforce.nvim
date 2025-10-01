package api

import (
	"fmt"
	"log/slog"
	"net/url"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func newLogMiddleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:   true,
		LogRemoteIP:  true,
		LogMethod:    true,
		LogURI:       true,
		LogRoutePath: true,
		LogUserAgent: true,
		LogStatus:    true,
		LogRequestID: true,
		LogError:     true,
		LogReferer:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			uri, _ := url.Parse(v.URI)
			q := uri.Query()
			if val := q.Get("id_token_hint"); val != "" {
				q.Set("id_token_hint", "filtered")
				uri.RawQuery = q.Encode()
			}

			attrs := []slog.Attr{
				slog.String("latencySec", fmt.Sprintf("%.3fs", v.Latency.Seconds())),
				slog.String("remoteIp", v.RemoteIP),
				slog.String("requestUrl", v.Method+" "+uri.String()),
				slog.String("requestRoute", v.RoutePath),
				slog.String("userAgent", v.UserAgent),
				slog.Int("responseCode", v.Status),
				slog.String("requestId", v.RequestID),
				slog.String("referer", v.Referer),
			}
			if v.Error != nil {
				attrs = append(attrs, slog.String("error", v.Error.Error()))
			}
			if userIDs, ok := c.Get(contextKeyUserIDs).(map[string]string); ok {
				attrs = append(attrs, slog.Any("currentUser", userIDs))
			}

			var l slog.Level
			if v.Status < 400 {
				l = slog.LevelInfo
			} else if v.Status < 500 {
				l = slog.LevelWarn
			} else {
				l = slog.LevelError
			}

			slog.LogAttrs(c.Request().Context(), l, "HTTP request", attrs...)
			return nil
		},
		Skipper: func(c echo.Context) bool {
			skippedPaths := []string{
				"/health",
				"/metrics",
			}

			return slices.ContainsFunc(skippedPaths, func(skippedPath string) bool {
				return skippedPath == c.Request().URL.Path
			})
		},
	})
}
