package repo

import (
	"backend/internal/domain/entity"
	"backend/pkg/utils"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepo abstracts user‑related persistence.
type UserRepo struct {
	pool *pgxpool.Pool
	rq   *entity.Queries
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool, rq: entity.New(pool)}
}

// CreateIfNotExists creates a user row when it does not exist.
func (r *UserRepo) CreateIfNotExists(ctx context.Context, id int64) error {
	if _, err := r.rq.GetUserById(ctx, id); err == nil {
		return nil
	}
	return utils.ExecInTx(ctx, r.pool, func(tq *entity.Queries) error { return tq.CreateUser(ctx, id) })
}

// Exists checks for a user record.
func (r *UserRepo) Exists(ctx context.Context, id int64) (bool, error) {
	_, err := r.rq.GetUserById(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" { // pgx.ErrNoRows unfetched, safe fallback
			return false, nil
		}
		return false, err
	}
	return true, nil
}
