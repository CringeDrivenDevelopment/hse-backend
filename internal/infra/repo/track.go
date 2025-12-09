package repo

import (
	"backend/internal/domain/entity"
	"backend/internal/transport/api/dto"
	"backend/pkg/utils"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TrackRepo struct {
	pool *pgxpool.Pool

	rq *entity.Queries
}

func NewTrackRepo(pool *pgxpool.Pool) *TrackRepo {
	return &TrackRepo{rq: entity.New(pool), pool: pool}
}

func (r *TrackRepo) AddBatch(ctx context.Context, tracks []dto.Track) error {
	return utils.ExecInTx(ctx, r.pool, func(tq *entity.Queries) error {
		for _, track := range tracks {
			_, err := r.rq.GetTrackById(ctx, track.Id)
			if errors.Is(err, pgx.ErrNoRows) {
				err = tq.CreateTrack(ctx, entity.CreateTrackParams{
					ID:        track.Id,
					Title:     track.Title,
					Authors:   track.Authors,
					Thumbnail: track.Thumbnail,
					Length:    track.Length,
					Explicit:  track.Explicit,
				})
				if err != nil {
					return err
				}
			}

			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *TrackRepo) GetById(ctx context.Context, id string) (entity.Track, error) {
	return r.rq.GetTrackById(ctx, id)
}

func (r *TrackRepo) Update(ctx context.Context, params entity.EditPlaylistParams) error {
	return utils.ExecInTx(ctx, r.pool, func(tq *entity.Queries) error {
		return tq.EditPlaylist(ctx, params)
	})
}
