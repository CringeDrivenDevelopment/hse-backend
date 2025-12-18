package services

import (
	"backend/internal/bot/dto"
	"backend/internal/domain/entity"
	"errors"

	"github.com/celestix/gotgproto/ext"
	"github.com/gotd/td/tg"
)

type ParticipantService struct{}

func NewParticipantService() *ParticipantService {
	return &ParticipantService{}
}

func (s *ParticipantService) HandleParticipant(update *ext.Update) (dto.Participant, error) {
	var data dto.Participant
	switch u := update.UpdateClass.(type) {
	case *tg.UpdateChannelParticipant:
		data = s.extractChannelData(u)
	case *tg.UpdateChatParticipant:
		data = s.extractChatData(u)
	default:
		return data, errors.New("invalid update type " + u.TypeName())
	}

	return data, nil
}

func (s *ParticipantService) extractChannelData(update *tg.UpdateChannelParticipant) dto.Participant {
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

	return dto.Participant{
		PrevRole: prevRole,
		NewRole:  newRole,
		UserID:   update.UserID,
		ChatID:   update.ChannelID,
		ActorID:  update.ActorID,
	}
}

func (s *ParticipantService) extractChatData(update *tg.UpdateChatParticipant) dto.Participant {
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

	return dto.Participant{
		PrevRole: prevRole,
		NewRole:  newRole,
		UserID:   update.UserID,
		ChatID:   update.ChatID,
		ActorID:  update.ActorID,
	}
}
