package repo

import (
	"backend/internal/api/dto"
	"backend/internal/domain/entity"
	"backend/pkg/utils"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PlaylistRepo encapsulates all database operations related to playlists.
type PlaylistRepo struct {
	pool *pgxpool.Pool
	rq   *entity.Queries
}

func NewPlaylistRepo(pool *pgxpool.Pool) *PlaylistRepo {
	return &PlaylistRepo{pool: pool, rq: entity.New(pool)}
}

// Create inserts a new playlist.
func (r *PlaylistRepo) Create(ctx context.Context, id, title string, ptype entity.PlaylistType, telegramId int64) error {
	return utils.ExecInTx(ctx, r.pool, func(tq *entity.Queries) error {
		return tq.CreatePlaylist(ctx, entity.CreatePlaylistParams{
			ID:            id,
			Title:         title,
			Thumbnail:     "",
			Tracks:        make([]string, 0),
			AllowedTracks: make([]string, 0),
			Type:          ptype,
			TelegramID:    telegramId,
		})
	})
}

// GetByGroup fetches the group's playlist.
func (r *PlaylistRepo) GetByGroup(ctx context.Context, telegramId int64) (entity.Playlist, error) {
	return r.rq.GetGroupPlaylist(ctx, telegramId)
}

// GetUserPlaylist fetches a playlist for a specific user.
func (r *PlaylistRepo) GetUserPlaylist(ctx context.Context, playlistId string, userId int64) (entity.GetUserPlaylistByIdRow, error) {
	return r.rq.GetUserPlaylistById(ctx, entity.GetUserPlaylistByIdParams{PlaylistID: playlistId, UserID: userId})
}

// GetAllByUser returns all playlists owned by a user.
func (r *PlaylistRepo) GetAllByUser(ctx context.Context, userId int64) ([]entity.GetUserPlaylistsRow, error) {
	return r.rq.GetUserPlaylists(ctx, userId)
}

// Edit updates mutable fields of a playlist.
func (r *PlaylistRepo) Edit(ctx context.Context, params entity.EditPlaylistParams) error {
	return utils.ExecInTx(ctx, r.pool, func(tq *entity.Queries) error { return tq.EditPlaylist(ctx, params) })
}

// Delete removes a playlist.
func (r *PlaylistRepo) Delete(ctx context.Context, playlistId string) error {
	return utils.ExecInTx(ctx, r.pool, func(tq *entity.Queries) error { return tq.DeletePlaylist(ctx, playlistId) })
}

// BuildDTO converts an entity.Playlist to the API DTO, fetching track details via TrackRepo.
func (r *PlaylistRepo) BuildDTO(ctx context.Context, pl entity.GetUserPlaylistByIdRow, tr *TrackRepo) (dto.Playlist, error) {
	// Resolve track DTOs
	tracks := make([]dto.Track, len(pl.Tracks))
	for i, tid := range pl.Tracks {
		t, err := tr.GetById(ctx, tid)
		if err != nil {
			return dto.Playlist{}, err
		}
		tracks[i] = dto.Track{Id: t.ID, Title: t.Title, Authors: t.Authors, Thumbnail: t.Thumbnail, Length: t.Length, Explicit: t.Explicit}
	}
	// Convert count fields safely
	count := 0
	if pl.Count.Valid {
		count = int(pl.Count.Int32)
	}
	allowedCount := 0
	if pl.AllowedCount.Valid {
		allowedCount = int(pl.AllowedCount.Int32)
	}
	length := 0
	return dto.Playlist{Id: pl.ID, Title: pl.Title, Thumbnail: pl.Thumbnail, Tracks: tracks, AllowedIds: pl.AllowedTracks, Count: count, Length: length, AllowedCount: allowedCount, Role: pl.Role, Type: string(pl.Type)}, nil
}

// BuildGroupDTO creates a DTO for a group playlist (entity.Playlist).
func (r *PlaylistRepo) BuildGroupDTO(ctx context.Context, pl entity.Playlist, tr *TrackRepo) (dto.Playlist, error) {
	// Resolve track DTOs
	tracks := make([]dto.Track, len(pl.Tracks))
	for i, tid := range pl.Tracks {
		t, err := tr.GetById(ctx, tid)
		if err != nil {
			return dto.Playlist{}, err
		}
		tracks[i] = dto.Track{Id: t.ID, Title: t.Title, Authors: t.Authors, Thumbnail: t.Thumbnail, Length: t.Length, Explicit: t.Explicit}
	}
	// Convert count fields safely
	count := 0
	if pl.Count.Valid {
		count = int(pl.Count.Int32)
	}
	allowedCount := 0
	if pl.AllowedCount.Valid {
		allowedCount = int(pl.AllowedCount.Int32)
	}
	length := 0
	// Return DTO (role and allowed ids are not relevant for group playlists)
	return dto.Playlist{Id: pl.ID, Title: pl.Title, Thumbnail: pl.Thumbnail, Tracks: tracks, AllowedIds: pl.AllowedTracks, Count: count, Length: length, AllowedCount: allowedCount, Role: "", Type: string(pl.Type)}, nil
}
