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

type Track struct {
	playlistService interfaces.PlaylistService
	trackService    interfaces.TrackService

	logger *zap.Logger
}

func NewTrack(playlistService *service.Playlist, trackService *service.Track, logger *zap.Logger, api huma.API, authMiddleware *middlewares.Auth) *Track {
	result := &Track{
		playlistService: playlistService,
		trackService:    trackService,
		logger:          logger,
	}

	result.setup(api, authMiddleware.IsAuthenticated)

	return result
}

// search - поиск трека по названию/исполнителю/...
func (h *Track) search(ctx context.Context, input *struct {
	Query string `query:"query"`
}) (*dto.SearchResponse, error) {
	val, ok := ctx.Value(middlewares.UserJwtKey).(int64)
	if !ok {
		err := utils.ErrContextUserNotFound

		return nil, utils.Convert(err, h.logger)
	}

	query := input.Query

	h.logger.Info(fmt.Sprintf("search: user_id - %d, query - %s", val, query))

	search, err := h.trackService.Search(ctx, "spotify", query)
	if err != nil {
		h.logger.Error(fmt.Sprintf("search error: user_id - %d, query - %s, error - %s", val, query, err.Error()))

		return nil, utils.Convert(err, h.logger)
	}

	return &dto.SearchResponse{Body: search}, nil
}

// submit - добавить трек в список на модерацию
func (h *Track) submit(ctx context.Context, input *dto.TrackAction) (*struct{}, error) {
	val, ok := ctx.Value(middlewares.UserJwtKey).(int64)
	if !ok {
		err := utils.ErrContextUserNotFound

		return nil, utils.Convert(err, h.logger)
	}

	h.logger.Info(fmt.Sprintf("submit: user_id - %d, playlist_id - %s, track_id - %s", val, input.PlaylistId, input.TrackId))

	if err := h.trackService.Submit(ctx, input.PlaylistId, input.TrackId, val); err != nil {
		h.logger.Error(fmt.Sprintf("submit error: user_id - %d, playlist_id - %s, track_id - %s, error - %s", val, input.PlaylistId, input.TrackId, err.Error()))

		return nil, utils.Convert(err, h.logger)
	}

	return nil, nil
}

// decline - удалить трек из списка на модерацию
func (h *Track) decline(ctx context.Context, input *dto.TrackAction) (*struct{}, error) {
	val, ok := ctx.Value(middlewares.UserJwtKey).(int64)
	if !ok {
		err := utils.ErrContextUserNotFound

		return nil, utils.Convert(err, h.logger)
	}

	h.logger.Info(fmt.Sprintf("decline: user_id - %d, playlist_id - %s, track_id - %s", val, input.PlaylistId, input.TrackId))

	if err := h.trackService.Decline(ctx, input.PlaylistId, input.TrackId, val); err != nil {
		h.logger.Error(fmt.Sprintf("decline error: user_id - %d, playlist_id - %s, track_id - %s, error - %s", val, input.PlaylistId, input.TrackId, err.Error()))

		return nil, utils.Convert(err, h.logger)
	}

	return nil, nil
}

// unapprove - удалить трек из списка разрешённых
func (h *Track) unapprove(ctx context.Context, input *dto.TrackAction) (*struct{}, error) {
	val, ok := ctx.Value(middlewares.UserJwtKey).(int64)
	if !ok {
		err := utils.ErrContextUserNotFound

		return nil, utils.Convert(err, h.logger)
	}

	h.logger.Info(fmt.Sprintf("unapprove: user_id - %d, playlist_id - %s, track_id - %s", val, input.PlaylistId, input.TrackId))

	if err := h.trackService.Unapprove(ctx, input.PlaylistId, input.TrackId, val); err != nil {
		h.logger.Error(fmt.Sprintf("unapprove error: user_id - %d, playlist_id - %s, track_id - %s, error - %s", val, input.PlaylistId, input.TrackId, err.Error()))

		return nil, utils.Convert(err, h.logger)
	}

	return nil, nil
}

// approve - добавить трек в список разрешённых
func (h *Track) approve(ctx context.Context, input *dto.TrackAction) (*struct{}, error) {
	val, ok := ctx.Value(middlewares.UserJwtKey).(int64)
	if !ok {
		err := utils.ErrContextUserNotFound

		return nil, utils.Convert(err, h.logger)
	}

	h.logger.Info(fmt.Sprintf("approve: user_id - %d, playlist_id - %s, track_id - %s", val, input.PlaylistId, input.TrackId))

	if err := h.trackService.Approve(ctx, input.PlaylistId, input.TrackId, val); err != nil {
		h.logger.Error(fmt.Sprintf("approve error: user_id - %d, playlist_id - %s, track_id - %s, error - %s", val, input.PlaylistId, input.TrackId, err.Error()))

		return nil, utils.Convert(err, h.logger)
	}

	return nil, nil
}
