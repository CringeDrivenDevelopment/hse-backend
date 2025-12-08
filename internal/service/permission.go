package service

import (
	"backend/internal/domain/entity"
	"backend/internal/transport/bot/models"
	"backend/pkg/utils"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Permission struct {
	pool *pgxpool.Pool
}

func NewPermission(pool *pgxpool.Pool) *Permission {
	return &Permission{pool: pool}
}

func (s *Permission) Add(ctx context.Context, role entity.PlaylistRole, playlist string, userId int64) error {
	rq := entity.New(s.pool)

	return utils.ExecInTx(ctx, s.pool, func(tq *entity.Queries) error {
		_, err := rq.GetUserById(ctx, userId)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}

			if err := tq.CreateUser(ctx, userId); err != nil {
				return err
			}
		}

		return tq.CreateRole(ctx, entity.CreateRoleParams{
			Role:       role,
			UserID:     userId,
			PlaylistID: playlist,
		})
	})
}

func (s *Permission) AddGroup(ctx context.Context, playlist string, users []models.ParticipantData) error {
	rq := entity.New(s.pool)

	return utils.ExecInTx(ctx, s.pool, func(tq *entity.Queries) error {
		for _, user := range users {
			_, err := rq.GetUserById(ctx, user.UserID)
			if err != nil {
				if !errors.Is(err, pgx.ErrNoRows) {
					return err
				}
				err = tq.CreateUser(ctx, user.UserID)
				if err != nil {
					return err
				}
			}

			err = tq.CreateRole(ctx, entity.CreateRoleParams{
				Role:       user.NewRole,
				UserID:     user.UserID,
				PlaylistID: playlist,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Permission) Remove(ctx context.Context, playlist string, userId int64) error {
	return utils.ExecInTx(ctx, s.pool, func(tq *entity.Queries) error {
		return tq.DeleteRole(ctx, entity.DeleteRoleParams{
			PlaylistID: playlist,
			UserID:     userId,
		})
	})
}

func (s *Permission) Edit(ctx context.Context, role entity.PlaylistRole, playlist string, userId int64) error {
	return utils.ExecInTx(ctx, s.pool, func(tq *entity.Queries) error {
		return tq.EditRole(ctx, entity.EditRoleParams{
			Role:       role,
			PlaylistID: playlist,
			UserID:     userId,
		})
	})
}

func (s *Permission) Get(ctx context.Context, userId int64, role entity.PlaylistRole) (string, error) {
	rq := entity.New(s.pool)

	return rq.GetRole(ctx, entity.GetRoleParams{
		Role:   role,
		UserID: userId,
	})
}
