package service

import (
	"backend/internal/infra/repo"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	repo *repo.UserRepo
}

func NewUser(pool *pgxpool.Pool) *User {
	return &User{repo: repo.NewUserRepo(pool)}
}

func (s *User) Create(ctx context.Context, id int64) error {
	return s.repo.CreateIfNotExists(ctx, id)
}

func (s *User) GetByID(ctx context.Context, id int64) error {
	_, err := s.repo.Exists(ctx, id)
	return err
}
