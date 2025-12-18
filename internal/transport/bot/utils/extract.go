package utils

import (
	"backend/internal/domain/entity"
	"backend/internal/transport/bot/models"

	"github.com/gotd/td/tg"
)

func ExtractChannelData(update *tg.UpdateChannelParticipant) models.ParticipantData {
	var newRole entity.PlaylistRole
	var prevRole entity.PlaylistRole

	if update.PrevParticipant != nil {
		switch update.PrevParticipant.(type) {
		case *tg.ChannelParticipant:
			prevRole = entity.PlaylistRoleViewer
		case *tg.ChannelParticipantAdmin:
			prevRole = entity.PlaylistRoleModerator
		case *tg.ChannelParticipantCreator:
			prevRole = entity.PlaylistRoleOwner
		case *tg.ChannelParticipantSelf:
			prevRole = entity.PlaylistRoleViewer
		default:
			prevRole = ""
		}
	}

	if update.NewParticipant != nil {
		switch update.NewParticipant.(type) {
		case *tg.ChannelParticipant:
			newRole = entity.PlaylistRoleViewer
		case *tg.ChannelParticipantAdmin:
			newRole = entity.PlaylistRoleModerator
		case *tg.ChannelParticipantCreator:
			newRole = entity.PlaylistRoleOwner
		case *tg.ChannelParticipantSelf:
			newRole = entity.PlaylistRoleViewer
		default:
			newRole = ""
		}
	}

	return models.ParticipantData{
		PrevRole: prevRole,
		NewRole:  newRole,
		UserID:   update.UserID,
		ChatID:   update.ChannelID,
		ActorID:  update.ActorID,
	}
}

func ExtractChatData(update *tg.UpdateChatParticipant) models.ParticipantData {
	var newRole entity.PlaylistRole
	var prevRole entity.PlaylistRole

	if update.PrevParticipant != nil {
		switch update.PrevParticipant.(type) {
		case *tg.ChatParticipant:
			prevRole = entity.PlaylistRoleViewer
		case *tg.ChatParticipantAdmin:
			prevRole = entity.PlaylistRoleModerator
		case *tg.ChatParticipantCreator:
			prevRole = entity.PlaylistRoleOwner
		default:
			prevRole = ""
		}
	}

	if update.NewParticipant != nil {
		switch update.NewParticipant.(type) {
		case *tg.ChatParticipant:
			newRole = entity.PlaylistRoleViewer
		case *tg.ChatParticipantAdmin:
			newRole = entity.PlaylistRoleModerator
		case *tg.ChatParticipantCreator:
			newRole = entity.PlaylistRoleOwner
		default:
			newRole = ""
		}
	}

	return models.ParticipantData{
		PrevRole: prevRole,
		NewRole:  newRole,
		UserID:   update.UserID,
		ChatID:   update.ChatID,
		ActorID:  update.ActorID,
	}
}
