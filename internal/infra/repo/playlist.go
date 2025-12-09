package repo

import (
	"backend/internal/domain/entity"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaylistRepo struct {
	pool *pgxpool.Pool

	rq *entity.Queries
}

func NewPlaylistRepo(pool *pgxpool.Pool) *PlaylistRepo {
	return &PlaylistRepo{rq: entity.New(pool), pool: pool}
}

func (r *PlaylistRepo) GetAll(ctx context.Context) ([]entity.Playlist, error) {
	return r.rq.GetPlaylist
}
