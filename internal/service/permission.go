package service

import (
	"backend/internal/bot/dto"
	"backend/internal/domain/entity"
	"backend/internal/infra/repo"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PermissionService struct {
	repo *repo.PermissionRepo
}

func NewPermissionService(pool *pgxpool.Pool) *PermissionService {
	return &PermissionService{repo: repo.NewPermissionRepo(pool)}
}

func (s *PermissionService) Add(ctx context.Context, role entity.PlaylistRole, playlist string, userId int64) error {
	return s.repo.Add(ctx, role, playlist, userId)
}

func (s *PermissionService) AddGroup(ctx context.Context, playlist string, users []dto.Participant) error {
	return s.repo.AddGroup(ctx, playlist, users)
}

func (s *PermissionService) Remove(ctx context.Context, playlist string, userId int64) error {
	return s.repo.Remove(ctx, playlist, userId)
}

func (s *PermissionService) Edit(ctx context.Context, role entity.PlaylistRole, playlist string, userId int64) error {
	return s.repo.Edit(ctx, role, playlist, userId)
}

func (s *PermissionService) Get(ctx context.Context, userId int64, role entity.PlaylistRole) (string, error) {
	return s.repo.Get(ctx, userId, role)
}
