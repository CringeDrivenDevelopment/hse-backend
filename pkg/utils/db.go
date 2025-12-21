package utils

import (
	"backend/internal/domain/entity"
	"context"
	"errors"
	"github.com/gotd/td/tg"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ExecInTx(ctx context.Context, pool *pgxpool.Pool, action func(tq *entity.Queries) error) error {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	tq := entity.New(tx)

	if err := action(tq); err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return txErr
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return txErr
		}
		return err
	}

	return nil
}

func GetChatID(peer tg.PeerClass) (int64, error) {
	switch p := peer.(type) {
	case *tg.PeerChannel:
		return p.ChannelID, nil
	case *tg.PeerChat:
		return p.ChatID, nil
	case *tg.PeerUser:
		return p.UserID, nil
	default:
		return 0, errors.New("invalid peer type")
	}
}
