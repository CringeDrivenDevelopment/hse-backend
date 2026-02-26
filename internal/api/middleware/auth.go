package middleware

import (
	"backend/internal/application/service"
	"backend/internal/interfaces"
	"context"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type Auth struct {
	authService interfaces.AuthService
	logger      *zap.Logger
}

const UserContextKey = "user"

func NewAuth(authService *service.AuthService, logger *zap.Logger) *Auth {

	return &Auth{
		authService: authService,
		logger:      logger,
	}
}

func (m *Auth) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		// Если пользователь не авторизован, то не будем ставить в контекст его данные, сами можем проверить в хендлере
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return next(c)
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, huma.Error401Unauthorized("invalid token type. only bearer is allowed"))
		}

		tokenString := strings.ReplaceAll(authHeader, "Bearer ", "")

		claims, err := m.authService.VerifyToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, huma.Error401Unauthorized("invalid or expired token"))
		}

		// Сохранить токен в контекст
		ctx := context.WithValue(c.Request().Context(), UserContextKey, claims)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
