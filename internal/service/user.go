package service

import (
	"backend/internal/infra/repo"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	repo *repo.UserRepo
}

func NewUserService(pool *pgxpool.Pool) *UserService {
	return &UserService{repo: repo.NewUserRepo(pool)}
}

func (s *UserService) Create(ctx context.Context, id int64) error {
	return s.repo.CreateIfNotExists(ctx, id)
}

func (s *UserService) GetByID(ctx context.Context, id int64) error {
	_, err := s.repo.Exists(ctx, id)
	return err
}
