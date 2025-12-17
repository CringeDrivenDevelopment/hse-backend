package handlers

import (
	"backend/internal/domain/entity"
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
	authService       interfaces.AuthService
	playlistService   interfaces.PlaylistService
	permissionService interfaces.PermissionService
	spotifyExportAPI  interfaces.SpotifyExportAPI

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

func (h *Playlist) export(ctx context.Context, input *struct {
	SpotifyData string `header:"X-Spotify-Data" doc:"Spotify Auth data for logging into user's account"`
	IsPublic    bool   `query:"is_public" doc:"Set to true if the playlist is public. Works only with spotify"`
	Id          string `path:"id" minLength:"26" maxLength:"26" example:"01JZ35PYGP6HJA08H0NHYPBHWD" doc:"playlist id"`
}) (*struct{}, error) {
	val, ok := ctx.Value(middlewares.UserJwtKey).(int64)
	if !ok {
		err := utils.ErrContextUserNotFound

		return nil, utils.Convert(err, h.logger)
	}

	var err error
	h.logger.Info(fmt.Sprintf("export: user_id - %d, playlist_id - %s", val, input.Id))

	pl, err := h.playlistService.GetById(ctx, input.Id, val)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("export error: user_id - %d, playlist_id - %s; error: %s", val, input.Id, err.Error()))

		return nil, utils.Convert(err, h.logger)
	}

	switch pl.Type {
	case string(entity.PlaylistTypeSpotify):
		// extract auth data, wtf, i can just use middleware for later, TODO: FIX THIS SHIT
		tok, err := h.authService.ParseSpotifyData(input.SpotifyData)
		if err != nil {
			h.logger.Warn(fmt.Sprintf("export error: user_id - %d; error: %s", val, err.Error()))

			return nil, utils.Convert(err, h.logger)
		}

		err = h.spotifyExportAPI.Export(ctx, tok, pl, input.IsPublic)
		if err != nil {
			h.logger.Warn(fmt.Sprintf("export error: user_id - %d; error: %s", val, err.Error()))

			return nil, utils.Convert(err, h.logger)
		}
	default:
		return nil, huma.Error400BadRequest("unsupported playlist type for export")
	}

	return nil, nil
}
