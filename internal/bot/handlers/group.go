package handlers

import (
	"backend/internal/bot/dto"
	botService "backend/internal/bot/service"
	"backend/internal/domain/entity"
	"backend/internal/interfaces"
	"backend/internal/service"
	"context"
	"errors"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/ext"
	"github.com/gotd/td/tg"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type GroupHandler struct {
	playlistService   interfaces.PlaylistService
	permissionService interfaces.PermissionService
	chatService       *botService.ChatService
	partService       *botService.ParticipantService
	client            *gotgproto.Client
	logger            *zap.Logger
}

func NewGroupHandler(
	playlistService *service.PlaylistService,
	permissionService *service.PermissionService,
	chatService *botService.ChatService,
	partService *botService.ParticipantService,
	client *gotgproto.Client,
	logger *zap.Logger,
) *GroupHandler {
	return &GroupHandler{
		playlistService:   playlistService,
		permissionService: permissionService,
		chatService:       chatService,
		partService:       partService,
		client:            client,
		logger:            logger,
	}
}

func (h *GroupHandler) HandleGroup(ctx *ext.Context, update *ext.Update) error {
	data, err := h.partService.HandleParticipant(update)
	if err != nil {
		h.logger.Error("failed to handle participant", zap.Error(err))
		return nil
	}

	h.logger.Info("handling group update", zap.Int64("chat_id", data.ChatID), zap.Int64("user_id", data.UserID), zap.String("prev_role", string(data.PrevRole)), zap.String("new_role", string(data.NewRole)))

	if data.UserID == h.client.Self.ID {
		return h.handleBot(ctx, &data)
	}
	return h.handleUser(ctx, &data)
}

func (h *GroupHandler) handleBot(ctx *ext.Context, data *dto.Participant) error {
	if data.PrevRole == "" {
		_, err := ctx.SendMessage(data.ActorID, &tg.MessagesSendMessageRequest{Message: "Привет, я Лотти - бот для управления плейлистами!" + "\n\n" +
			"Можешь дать мне права администратора (чтобы я мог видеть лог AKA recent actions/недавние действия, а также всех участников и администраторов)"})
		if err != nil {
			h.logger.Error("failed to send message", zap.Error(err))
		}
		return nil
	}

	if data.PrevRole == entity.PlaylistRoleViewer && data.NewRole == entity.PlaylistRoleModerator {
		_, err := ctx.SendMessage(data.ActorID, &tg.MessagesSendMessageRequest{Message: "Респект тебе за админку, сейчас создам плейлист!"})
		if err != nil {
			h.logger.Error("failed to send message", zap.Error(err))
			return nil
		}

		chat, err := h.chatService.GetChatInfo(ctx.Context, data.ChatID, data.ActorID)
		if err != nil {
			h.logger.Error("failed to get chat info", zap.Error(err))
			return nil
		}

		create, err := h.playlistService.Create(ctx.Context, chat.Title, entity.PlaylistTypeYoutube, data.ChatID)
		if err != nil {
			h.logger.Error("failed to create playlist", zap.Error(err))
			return nil
		}

		err = h.permissionService.AddGroup(ctx.Context, create.Id, *chat.Users)
		if err != nil {
			h.logger.Error("failed to add group", zap.Error(err))
			return nil
		}

		return nil
	}

	if data.NewRole == "" || data.NewRole == entity.PlaylistRoleViewer {
		playlist, err := h.playlistService.GetByGroup(ctx, data.ChatID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			h.logger.Error("failed to get playlist by group", zap.Error(err))
			return nil
		}

		err = h.playlistService.Delete(ctx.Context, playlist.Id)
		if err != nil {
			h.logger.Error("failed to delete playlist", zap.Error(err))
			return nil
		}

		return nil
	}
	return nil
}

func (h *GroupHandler) handleUser(ctx context.Context, data *dto.Participant) error {
	playlist, err := h.playlistService.GetByGroup(ctx, data.ChatID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		h.logger.Error("failed to get playlist by group", zap.Error(err))
		return nil
	}

	if data.NewRole == "" {
		err = h.permissionService.Remove(ctx, playlist.Id, data.UserID)
		if err != nil {
			h.logger.Error("failed to remove permission", zap.Error(err))
			return nil
		}
	} else {
		if data.PrevRole == "" {
			err = h.permissionService.Add(ctx, data.NewRole, playlist.Id, data.UserID)
			if err != nil {
				h.logger.Error("failed to add permission", zap.Error(err))
				return nil
			}
		} else {
			err = h.permissionService.Edit(ctx, data.NewRole, playlist.Id, data.UserID)
			if err != nil {
				h.logger.Error("failed to edit permission", zap.Error(err))
				return nil
			}
		}
	}
	return nil
}
