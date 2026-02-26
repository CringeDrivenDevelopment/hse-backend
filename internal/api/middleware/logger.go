package middleware

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.uber.org/zap"
)

func NewLogger(log *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:       true,
		LogURI:          true,
		LogMethod:       true,
		LogHost:         true,
		LogRemoteIP:     true,
		LogUserAgent:    true,
		LogResponseSize: true,
		LogRequestID:    true,
		LogLatency:      true,
		HandleError:     true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			fields := []zap.Field{
				zap.String("remote_ip", v.RemoteIP),
				zap.String("latency", v.Latency.String()),
				zap.String("host", v.Host),
				zap.String("request", fmt.Sprintf("%s %s", v.Method, v.URI)),
				zap.Int("status", v.Status),
				zap.Int64("size", v.ResponseSize),
				zap.String("user_agent", v.UserAgent),
				zap.String("request_id", v.RequestID),
			}
			n := v.Status
			switch {
			case n >= 500:
				log.With(zap.Error(v.Error)).Error("Server error", fields...)
			case n >= 400:
				log.Warn("Client error", fields...)
			case n >= 300:
				log.Info("Redirection", fields...)
			default:
				log.Info("Success", fields...)
			}

			return nil
		},
	})
}
