package infra

import (
	"context"
	"errors"
	"net/http"
	"time"

	appMiddleware "backend/internal/api/middleware"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewEcho(lc fx.Lifecycle, logger *zap.Logger, authMiddleware *appMiddleware.Auth) *echo.Echo {
	router := echo.New()

	router.Use(echoMiddleware.Recover())
	router.Use(echoMiddleware.RequestID())
	router.Use(appMiddleware.NewLogger(logger))
	router.Use(authMiddleware.Handle)
	router.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{
			"https://tg-mini-app.local",
			"https://muse-ghp.lxft.tech",
			"https://muse.lxft.tech",
			"https://localhost",
			"http://localhost",
		},
	}))

	router.GET("/api/ping", func(c *echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	shutdownCh := make(chan error, 1)
	serverCtx, serverCancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				sc := echo.StartConfig{
					Address:         ":8080",
					GracefulTimeout: 5 * time.Second,
					HideBanner:      true,
					HidePort:        true,
				}
				logger.Info("api server started on :8080")

				if err := sc.Start(serverCtx, router); err != nil && !errors.Is(err, http.ErrServerClosed) && !errors.Is(err, context.Canceled) {
					logger.Error("server error", zap.Error(err))
					shutdownCh <- err
				} else {
					shutdownCh <- nil
				}
				close(shutdownCh)
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("stopping api server")
			serverCancel()

			select {
			case err := <-shutdownCh:
				if err != nil {
					logger.Error("api server shutdown with error", zap.Error(err))
					return err
				}
				logger.Info("api server stopped")
			case <-ctx.Done():
				logger.Error("api server shutdown timeout")
				return ctx.Err()
			}

			return nil
		},
	})

	return router
}
