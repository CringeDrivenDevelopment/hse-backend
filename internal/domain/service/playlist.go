package service

import (
	"backend/internal/api/dto"
	"backend/internal/domain/entity"
	"backend/internal/infra/repo"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

type PlaylistService struct {
	playlistRepo *repo.PlaylistRepo
	trackRepo    *repo.TrackRepo
}

// TODO: fix constructor
func NewPlaylistService(pool *pgxpool.Pool, trackRepo *repo.TrackRepo) *PlaylistService {
	return &PlaylistService{playlistRepo: repo.NewPlaylistRepo(pool), trackRepo: trackRepo}
}

func (s *PlaylistService) Create(ctx context.Context, title string, playlistType entity.PlaylistType, telegramId int64) (dto.Playlist, error) {
	id := ulid.Make().String()
	if err := s.playlistRepo.Create(ctx, id, title, playlistType, telegramId); err != nil {
		return dto.Playlist{}, err
	}
	return dto.Playlist{Id: id, Title: title, Type: string(playlistType)}, nil
}

func (s *PlaylistService) GetByGroup(ctx context.Context, telegramId int64) (dto.Playlist, error) {
	pl, err := s.playlistRepo.GetByGroup(ctx, telegramId)
	if err != nil {
		return dto.Playlist{}, err
	}
	return s.playlistRepo.BuildGroupDTO(ctx, pl, s.trackRepo)
}

func (s *PlaylistService) GetById(ctx context.Context, playlistId string, userId int64) (dto.Playlist, error) {
	pl, err := s.playlistRepo.GetUserPlaylist(ctx, playlistId, userId)
	if err != nil {
		return dto.Playlist{}, err
	}
	return s.playlistRepo.BuildDTO(ctx, pl, s.trackRepo)
}

func (s *PlaylistService) GetAll(ctx context.Context, userId int64) ([]dto.Playlist, error) {
	playlists, err := s.playlistRepo.GetAllByUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	result := make([]dto.Playlist, len(playlists))
	for i, pl := range playlists {
		count, allowedCount, length := 0, 0, int(pl.Time)
		if pl.Count.Valid {
			count = int(pl.Count.Int32)
		}
		if pl.AllowedCount.Valid {
			allowedCount = int(pl.AllowedCount.Int32)
		}
		result[i] = dto.Playlist{Id: pl.ID, Title: pl.Title, Thumbnail: pl.Thumbnail, Count: count, Length: length, AllowedCount: allowedCount, Tracks: []dto.Track{}, AllowedIds: []string{}, Role: pl.Role, Type: string(pl.Type)}
	}
	return result, nil
}

func (s *PlaylistService) Rename(ctx context.Context, playlistId, title string, userId int64) error {
	pl, err := s.playlistRepo.GetUserPlaylist(ctx, playlistId, userId)
	if err != nil {
		return err
	}
	pl.Title = title
	return s.playlistRepo.Edit(ctx, entity.EditPlaylistParams{ID: pl.ID, Title: pl.Title})
}

func (s *PlaylistService) UpdatePhoto(ctx context.Context, playlistId, thumbnail string, userId int64) error {
	pl, err := s.playlistRepo.GetUserPlaylist(ctx, playlistId, userId)
	if err != nil {
		return err
	}
	pl.Thumbnail = thumbnail
	return s.playlistRepo.Edit(ctx, entity.EditPlaylistParams{ID: pl.ID, Thumbnail: pl.Thumbnail})
}

func (s *PlaylistService) Delete(ctx context.Context, playlistId string) error {
	// Delegate deletion to the repository layer to keep DB logic out of the service.
	return s.playlistRepo.Delete(ctx, playlistId)
}
