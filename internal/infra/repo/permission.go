package repo

import (
    "backend/internal/domain/entity"
    "backend/internal/transport/bot/models"
    "backend/pkg/utils"
    "context"
    "errors"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

// PermissionRepo provides database operations for permission/role management.
type PermissionRepo struct {
    pool *pgxpool.Pool
    rq   *entity.Queries
}

func NewPermissionRepo(pool *pgxpool.Pool) *PermissionRepo {
    return &PermissionRepo{pool: pool, rq: entity.New(pool)}
}

func (r *PermissionRepo) Add(ctx context.Context, role entity.PlaylistRole, playlist string, userId int64) error {
    return utils.ExecInTx(ctx, r.pool, func(tq *entity.Queries) error {
        // Ensure user exists
        if _, err := r.rq.GetUserById(ctx, userId); err != nil {
            if !errors.Is(err, pgx.ErrNoRows) {
                return err
            }
            if err := tq.CreateUser(ctx, userId); err != nil {
                return err
            }
        }
        // Create role mapping
        return tq.CreateRole(ctx, entity.CreateRoleParams{Role: role, UserID: userId, PlaylistID: playlist})
    })
}

func (r *PermissionRepo) AddGroup(ctx context.Context, playlist string, users []models.ParticipantData) error {
    return utils.ExecInTx(ctx, r.pool, func(tq *entity.Queries) error {
        for _, user := range users {
            if _, err := r.rq.GetUserById(ctx, user.UserID); err != nil {
                if !errors.Is(err, pgx.ErrNoRows) {
                    return err
                }
                if err = tq.CreateUser(ctx, user.UserID); err != nil {
                    return err
                }
            }
            if err := tq.CreateRole(ctx, entity.CreateRoleParams{Role: user.NewRole, UserID: user.UserID, PlaylistID: playlist}); err != nil {
                return err
            }
        }
        return nil
    })
}

func (r *PermissionRepo) Remove(ctx context.Context, playlist string, userId int64) error {
    return utils.ExecInTx(ctx, r.pool, func(tq *entity.Queries) error {
        return tq.DeleteRole(ctx, entity.DeleteRoleParams{PlaylistID: playlist, UserID: userId})
    })
}

func (r *PermissionRepo) Edit(ctx context.Context, role entity.PlaylistRole, playlist string, userId int64) error {
    return utils.ExecInTx(ctx, r.pool, func(tq *entity.Queries) error {
        return tq.EditRole(ctx, entity.EditRoleParams{Role: role, PlaylistID: playlist, UserID: userId})
    })
}

func (r *PermissionRepo) Get(ctx context.Context, userId int64, role entity.PlaylistRole) (string, error) {
    return r.rq.GetRole(ctx, entity.GetRoleParams{Role: role, UserID: userId})
}

