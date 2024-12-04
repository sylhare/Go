package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var logger, _ = zap.NewProduction()

func NewLogValuesFunc(logger *zap.Logger) func(c echo.Context, v middleware.RequestLoggerValues) error {
	return func(c echo.Context, v middleware.RequestLoggerValues) error {
		fields := []zap.Field{
			zap.String("remote_ip", v.RemoteIP),
			zap.String("host", v.Host),
			zap.String("uri", v.URI),
			zap.String("method", v.Method),
			zap.Int("status", v.Status),
			zap.Duration("latency", v.Latency),
			zap.String("latency_human", v.Latency.String()),
			zap.String("user_agent", v.UserAgent),
		}

		if v.Error != nil {
			logger.Error("request failed", append(fields, zap.Error(v.Error))...)
		} else {
			logger.Info("request completed", fields...)
		}

		return nil
	}
}

var zapLoggerConfig = middleware.RequestLoggerConfig{
	LogRemoteIP:   true,
	LogHost:       true,
	LogURI:        true,
	LogMethod:     true,
	LogStatus:     true,
	LogLatency:    true,
	LogUserAgent:  true,
	LogError:      true,
	LogValuesFunc: NewLogValuesFunc(logger),
}

var Logger = middleware.RequestLoggerWithConfig(zapLoggerConfig)
