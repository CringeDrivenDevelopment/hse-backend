package service

import (
	"backend/internal/domain/entity"
	"backend/internal/infra/repo"
	"backend/internal/transport/bot/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Permission struct {
	repo *repo.PermissionRepo
}

func NewPermission(pool *pgxpool.Pool) *Permission {
	return &Permission{repo: repo.NewPermissionRepo(pool)}
}

func (s *Permission) Add(ctx context.Context, role entity.PlaylistRole, playlist string, userId int64) error {
	return s.repo.Add(ctx, role, playlist, userId)
}

func (s *Permission) AddGroup(ctx context.Context, playlist string, users []models.ParticipantData) error {
	return s.repo.AddGroup(ctx, playlist, users)
}

func (s *Permission) Remove(ctx context.Context, playlist string, userId int64) error {
	return s.repo.Remove(ctx, playlist, userId)
}

func (s *Permission) Edit(ctx context.Context, role entity.PlaylistRole, playlist string, userId int64) error {
	return s.repo.Edit(ctx, role, playlist, userId)
}

func (s *Permission) Get(ctx context.Context, userId int64, role entity.PlaylistRole) (string, error) {
	return s.repo.Get(ctx, userId, role)
}
