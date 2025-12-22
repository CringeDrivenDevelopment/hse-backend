package service

import (
	"backend/internal/domain/entity"
	"backend/internal/infra/repo"
	"backend/internal/interfaces"
	"backend/internal/transport/api/dto"
	"backend/pkg/spotify"
	"backend/pkg/utils"
	"backend/pkg/youtube"
	"context"
	"slices"

	"github.com/jackc/pgx/v5"
)

type TrackService struct {
	trackRepo    *repo.TrackRepo
	playlistRepo *repo.PlaylistRepo

	youtube interfaces.SearchAPI
	spotify interfaces.SearchAPI
}

func NewTrackService(trackRepo *repo.TrackRepo, playlistRepo *repo.PlaylistRepo, ytApi *youtube.API, spotify *spotify.API) *TrackService {
	return &TrackService{trackRepo: trackRepo, playlistRepo: playlistRepo, youtube: ytApi, spotify: spotify}
}

/*
Search - метод для поиска треков на какой-либо из площадок

Старая сигнатура:
Search(ctx context.Context, query string, userId int64) ([]dto.Track, error)
Изменилась, из-за поддержки множества площадок, а так же, из-за ненадобности получения списка плейлистов для трека
*/
func (s *TrackService) Search(ctx context.Context, platform, query string) ([]dto.Track, error) {
	var tracks []dto.Track
	var err error

	switch platform {
	case string(entity.PlaylistTypeYoutube):
		tracks, err = s.youtube.Search(ctx, query)
	case string(entity.PlaylistTypeSpotify):
		tracks, err = s.spotify.Search(ctx, query)
	default:
		err = utils.ErrUnknownPlatform
	}

	if err != nil {
		return nil, err
	}

	if err := s.trackRepo.AddBatch(ctx, tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (s *TrackService) GetById(ctx context.Context, id string) (dto.Track, error) {
	track, err := s.trackRepo.GetById(ctx, id)
	if err != nil {
		return dto.Track{}, err
	}

	return dto.Track{
		Id:        track.ID,
		Title:     track.Title,
		Authors:   track.Authors,
		Thumbnail: track.Thumbnail,
		Length:    track.Length,
		Explicit:  track.Explicit,
	}, nil
}

func (s *TrackService) Approve(ctx context.Context, playlistId, trackId string, userId int64) error {
	playlist, err := s.playlistRepo.GetUserPlaylist(ctx, playlistId, userId)
	if err != nil {
		return err
	}

	if playlist.Role != entity.PlaylistRoleOwner && playlist.Role != entity.PlaylistRoleModerator {
		return utils.ErrNotEnoughPerms
	}

	if slices.Contains(playlist.AllowedTracks, trackId) {
		return nil
	}

	if !slices.Contains(playlist.Tracks, trackId) {
		return pgx.ErrNoRows
	}

	playlist.AllowedTracks = append(playlist.AllowedTracks, trackId)

	if err := s.trackRepo.Update(ctx, entity.EditPlaylistParams{
		ID:            playlistId,
		AllowedTracks: playlist.AllowedTracks,
		Type:          playlist.Type,
	}); err != nil {
		return err
	}

	return nil
}

func (s *TrackService) Decline(ctx context.Context, playlistId, trackId string, userId int64) error {
	playlist, err := s.playlistRepo.GetUserPlaylist(ctx, playlistId, userId)
	if err != nil {
		return err
	}

	if playlist.Role != entity.PlaylistRoleOwner && playlist.Role != entity.PlaylistRoleModerator {
		return utils.ErrNotEnoughPerms
	}

	if !slices.Contains(playlist.Tracks, trackId) {
		return pgx.ErrNoRows
	}

	for i, track := range playlist.Tracks {
		if track == trackId {
			playlist.Tracks = append(playlist.Tracks[:i], playlist.Tracks[i+1:]...)
			break
		}
	}

	if err := s.trackRepo.Update(ctx, entity.EditPlaylistParams{
		ID:            playlistId,
		AllowedTracks: playlist.AllowedTracks,
		Type:          playlist.Type,
	}); err != nil {
		return err
	}

	return nil
}

func (s *TrackService) Submit(ctx context.Context, playlistId, trackId string, userId int64) error {
	playlist, err := s.playlistRepo.GetUserPlaylist(ctx, playlistId, userId)
	if err != nil {
		return err
	}

	if _, err := s.trackRepo.GetById(ctx, trackId); err != nil {
		return err
	}

	tracks := playlist.Tracks
	allowedTracks := playlist.AllowedTracks
	if (playlist.Role == entity.PlaylistRoleOwner || playlist.Role == entity.PlaylistRoleModerator) && !slices.Contains(allowedTracks, trackId) {
		tracks = append(tracks, trackId)
		allowedTracks = append(allowedTracks, trackId)
	} else if !slices.Contains(tracks, trackId) && playlist.Role == entity.PlaylistRoleViewer {
		tracks = append(tracks, trackId)
	} else {
		return nil
	}

	if err := s.trackRepo.Update(ctx, entity.EditPlaylistParams{
		ID:            playlistId,
		Tracks:        tracks,
		AllowedTracks: allowedTracks,
		Type:          playlist.Type,
	}); err != nil {
		return err
	}

	return nil
}

func (s *TrackService) Unapprove(ctx context.Context, playlistId, trackId string, userId int64) error {
	playlist, err := s.playlistRepo.GetUserPlaylist(ctx, playlistId, userId)
	if err != nil {
		return err
	}

	if playlist.Role != entity.PlaylistRoleOwner && playlist.Role != entity.PlaylistRoleModerator {
		return utils.ErrNotEnoughPerms
	}

	if !slices.Contains(playlist.AllowedTracks, trackId) && !slices.Contains(playlist.Tracks, trackId) {
		return pgx.ErrNoRows
	}

	for i, track := range playlist.AllowedTracks {
		if track == trackId {
			playlist.AllowedTracks = append(playlist.AllowedTracks[:i], playlist.AllowedTracks[i+1:]...)
			break
		}
	}

	if err := s.trackRepo.Update(ctx, entity.EditPlaylistParams{
		ID:            playlistId,
		AllowedTracks: playlist.AllowedTracks,
		Type:          playlist.Type,
	}); err != nil {
		return err
	}

	return nil
}
