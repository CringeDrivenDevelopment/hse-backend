package handlers

import (
	"backend/internal/interfaces"
	"backend/internal/service"
	"backend/internal/transport/api/dto"
	"backend/internal/transport/api/middlewares"
	"backend/pkg/utils"
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

type Playlist struct {
	playlistService   interfaces.PlaylistService
	permissionService interfaces.PermissionService

	logger *zap.Logger
}

// NewPlaylist - создать новый экземпляр обработчика
func NewPlaylist(playlistService *service.Playlist, permissionService *service.Permission, logger *zap.Logger, api huma.API, authMiddleware *middlewares.Auth) *Playlist {
	result := &Playlist{
		playlistService:   playlistService,
		permissionService: permissionService,
		logger:            logger,
	}

	result.setup(api, authMiddleware.IsAuthenticated)

	return result
}

// getById - получить плейлист по ID
func (h *Playlist) getById(ctx context.Context, input *struct {
	Id string `path:"id" minLength:"26" maxLength:"26" example:"01JZ35PYGP6HJA08H0NHYPBHWD" doc:"playlist id"`
}) (*dto.PlaylistByIdResponse, error) {
	val, ok := ctx.Value(middlewares.UserJwtKey).(int64)
	if !ok {
		err := utils.ErrContextUserNotFound

		return nil, utils.Convert(err, h.logger)
	}

	h.logger.Info(fmt.Sprintf("playlistById: user_id - %d, playlist_id - %s", val, input.Id))

	resp, err := h.playlistService.GetById(ctx, input.Id, val)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("playlistById error: user_id - %d, playlist_id - %s; error: %s", val, input.Id, err.Error()))

		return nil, utils.Convert(err, h.logger)
	}

	return &dto.PlaylistByIdResponse{Body: resp}, nil
}

// getAll - получить список плейлистов для пользователя
func (h *Playlist) getAll(ctx context.Context, _ *struct{}) (*dto.PlaylistsResponse, error) {
	val, ok := ctx.Value(middlewares.UserJwtKey).(int64)
	if !ok {
		err := utils.ErrContextUserNotFound

		return nil, utils.Convert(err, h.logger)
	}

	h.logger.Info(fmt.Sprintf("playlists: user_id - %d", val))

	resp, err := h.playlistService.GetAll(ctx, val)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("playlists error: user_id - %d; error: %s", val, err.Error()))

		return nil, utils.Convert(err, h.logger)
	}

	return &dto.PlaylistsResponse{Body: resp}, nil
}
