package handlers

import (
	"backend/internal/interfaces"
	"backend/internal/service"
	"backend/internal/transport/api/dto"
	"backend/pkg/utils"
	"context"
	"errors"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Auth struct {
	userService interfaces.UserService
	authService interfaces.AuthService

	logger *zap.Logger
}

// NewAuth - создать новый экземпляр обработчика
func NewAuth(userService *service.UserService, authService *service.AuthService, logger *zap.Logger, api huma.API) *Auth {
	result := &Auth{
		userService: userService,
		authService: authService,
		logger:      logger,
	}

	result.setup(api)

	return result
}

// login - Получить токен для взаимодействия. Нуждается в Raw строке из Telegram Mini App. Действует 1 час
func (h *Auth) login(ctx context.Context, input *dto.AuthInputStruct) (*dto.AuthOutputStruct, error) {
	h.logger.Info("login: " + input.Body.Raw)

	id, err := h.authService.ParseTelegramData(input.Body.Raw)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("login error: initdata - %s, error - %s", input.Body.Raw, err.Error()))

		return nil, utils.Convert(err, h.logger)
	}

	if err := h.userService.GetByID(ctx, id); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			h.logger.Warn(fmt.Sprintf("login error: initdata - %s, error - %s", input.Body.Raw, err.Error()))

			return nil, utils.Convert(err, h.logger)
		}

		if err := h.userService.Create(ctx, id); err != nil {
			h.logger.Warn(fmt.Sprintf("login error: initdata - %s, error - %s", input.Body.Raw, err.Error()))

			return nil, utils.Convert(err, h.logger)
		}
	}

	token, err := h.authService.GenerateToken(id)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("login error: initdata - %s, error - %s", input.Body.Raw, err.Error()))

		return nil, utils.Convert(err, h.logger)
	}

	tokenData := dto.Token{
		Token: token,
	}

	return &dto.AuthOutputStruct{Body: tokenData}, nil
}
