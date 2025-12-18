package handlers

import (
	"context"
	"errors"

	"backend/internal/interfaces"

	"github.com/celestix/gotgproto/ext"
	"github.com/gotd/td/tg"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type ChatActionHandler struct {
	playlistService interfaces.PlaylistService
	logger          *zap.Logger
}

func NewChatActionHandler(playlistService interfaces.PlaylistService, logger *zap.Logger) *ChatActionHandler {
	return &ChatActionHandler{
		playlistService: playlistService,
		logger:          logger,
	}
}

func (h *ChatActionHandler) HandleChatAction(ctx *ext.Context, update *ext.Update) error {
	var msg tg.MessageClass
	switch um := update.UpdateClass.(type) {
	case *tg.UpdateNewMessage:
		msg = um.Message
	case *tg.UpdateNewChannelMessage:
		msg = um.Message
	default:
		return nil
	}

	serviceMessage, ok := msg.(*tg.MessageService)
	if !ok {
		return nil
	}

	chatID, err := ctx.GetChatID()
	if err != nil {
		h.logger.Error("failed to get chat id", zap.Error(err))
		return nil
	}

	h.logger.Info("handle service message", zap.String("type", serviceMessage.TypeName()), zap.Int("id", serviceMessage.ID), zap.Int64("chat_id", chatID))

	switch smResult := serviceMessage.Action.(type) {
	case *tg.MessageActionChatEditTitle:
		err := h.handleTitleUpdate(ctx.Context, smResult.Title, chatID)
		if err != nil {
			h.logger.Error("failed to handle title update", zap.Error(err))
		}
	}

	return nil
}

func (h *ChatActionHandler) handleTitleUpdate(ctx context.Context, title string, chatID int64) error {
	playlist, err := h.playlistService.GetByGroup(ctx, chatID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	return h.playlistService.Rename(ctx, playlist.Id, title, chatID)
}
