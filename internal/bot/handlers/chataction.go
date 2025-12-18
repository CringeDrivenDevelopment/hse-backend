package handlers

import (
	"context"
	"errors"

	"backend/internal/interfaces"
	"backend/pkg/utils"

	"github.com/celestix/gotgproto/ext"
	"github.com/gotd/td/tg"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type ChatActionHandler struct {
	playlistService    interfaces.PlaylistService
	participantService interfaces.ParticipantService
	logger             *zap.Logger
}

func NewChatActionHandler(
	playlistService interfaces.PlaylistService,
	participantService interfaces.ParticipantService,
	logger *zap.Logger,
) *ChatActionHandler {
	return &ChatActionHandler{
		playlistService:    playlistService,
		participantService: participantService,
		logger:             logger,
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

	chatID, err := utils.GetChatID(serviceMessage.PeerID)
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
	case *tg.MessageActionChatAddUser:
		err = h.participantService.Add(ctx.Context, chatID, smResult.Users)
		if err != nil {
			h.logger.Error("failed to handle add user", zap.Error(err))
		}
	case *tg.MessageActionChatDeleteUser:
		err = h.participantService.Remove(ctx.Context, chatID, smResult.UserID)
		if err != nil {
			h.logger.Error("failed to handle delete user", zap.Error(err))
		}
	case *tg.MessageActionChatJoinedByLink:
		fromID, ok := serviceMessage.FromID.(*tg.PeerUser)
		if !ok {
			return errors.New("failed to handle joined by link")
		}

		err = h.participantService.Add(ctx.Context, chatID, []int64{fromID.UserID})
		if err != nil {
			h.logger.Error("failed to handle join by link user", zap.Error(err))
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
