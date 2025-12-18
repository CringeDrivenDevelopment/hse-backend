package services

import (
	"backend/internal/bot/dto"
	"backend/internal/domain/entity"
	"context"
	"errors"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/functions"
	"github.com/gotd/td/tg"
)

type ChatService struct {
	client *gotgproto.Client
}

func NewChatService(client *gotgproto.Client) *ChatService {
	return &ChatService{client: client}
}

func (s *ChatService) GetChatInfo(ctx context.Context, chatID, actorID int64) (*dto.ChatInfo, error) {
	peer := functions.GetInputPeerClassFromId(s.client.PeerStorage, chatID)
	if peer == nil {
		return nil, errors.New("peer is nil")
	}

	var chat tg.ChatClass
	var info *dto.ChatInfo
	var users *[]dto.Participant
	var err error

	switch peerResult := peer.(type) {
	case *tg.InputPeerChat:
		chatFull, err := s.client.API().MessagesGetFullChat(ctx, chatID)
		if err != nil {
			return nil, err
		}

		if len(chatFull.Chats) != 1 {
			return nil, errors.New("no chat found")
		}

		chat = chatFull.Chats[0]
		var usersTemp []dto.Participant
		for _, user := range chatFull.Users {
			if val, ok := user.(*tg.User); ok {
				if val.Bot {
					continue
				}

				if val.ID != actorID {
					usersTemp = append(usersTemp, dto.Participant{
						NewRole: entity.PlaylistRoleViewer,
						ChatID:  chatID,
						UserID:  val.ID,
					})
				} else {
					usersTemp = append(usersTemp, dto.Participant{
						NewRole: entity.PlaylistRoleOwner,
						ChatID:  chatID,
						UserID:  val.ID,
					})
				}
			}
		}
		users = &usersTemp
	case *tg.InputPeerChannel:
		channelFull, err := s.client.API().ChannelsGetFullChannel(ctx, &tg.InputChannel{
			ChannelID:  peerResult.ChannelID,
			AccessHash: peerResult.AccessHash,
		})
		if err != nil {
			return nil, err
		}

		if len(channelFull.Chats) != 1 {
			return nil, errors.New("no chat found")
		}

		chat = channelFull.Chats[0]
	default:
		return nil, errors.New("unknown peer type " + peerResult.TypeName())
	}

	if chat == nil {
		return nil, errors.New("no chat found")
	}

	switch c := chat.(type) {
	case *tg.Chat:
		info = &dto.ChatInfo{
			Title: c.Title,
			Photo: c.Photo,
		}
		if users == nil || len(*users) == 0 {
			return info, errors.New("no chat users found")
		}
	case *tg.Channel:
		info = &dto.ChatInfo{
			Title: c.Title,
			Photo: c.Photo,
		}

		users, err = s.iterateParticipants(ctx, &tg.InputChannel{
			AccessHash: c.AccessHash,
			ChannelID:  c.ID,
		})

		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unknown channel type " + c.TypeName())
	}

	info.Users = users

	return info, nil
}

func (s *ChatService) iterateParticipants(ctx context.Context, channel *tg.InputChannel) (*[]dto.Participant, error) {
	seen := make(map[int64]bool)
	var data []dto.Participant

	filters := []tg.ChannelParticipantsFilterClass{
		&tg.ChannelParticipantsRecent{},
		&tg.ChannelParticipantsAdmins{},
	}

	for _, filter := range filters {
		offset := 0
		const limit = 100

		for {
			resp, err := s.client.API().ChannelsGetParticipants(ctx, &tg.ChannelsGetParticipantsRequest{
				Channel: channel,
				Filter:  filter,
				Offset:  offset,
				Limit:   limit,
				Hash:    0,
			})
			if err != nil {
				return nil, err
			}

			val, ok := resp.(*tg.ChannelsChannelParticipants)
			if !ok {
				return nil, errors.New("invalid response " + resp.TypeName())
			}

			if len(val.Participants) == 0 {
				break
			}

			for _, participant := range val.Participants {
				userId, role := s.extractUserAndRole(participant)
				if userId == 0 || userId == s.client.Self.ID || seen[userId] {
					continue
				}

				seen[userId] = true
				data = append(data, dto.Participant{
					UserID:  userId,
					NewRole: role,
					ChatID:  channel.ChannelID,
				})
			}

			offset += len(val.Participants)
		}
	}

	return &data, nil
}

func (s *ChatService) extractUserAndRole(participant tg.ChannelParticipantClass) (int64, entity.PlaylistRole) {
	switch p := participant.(type) {
	case *tg.ChannelParticipant:
		return p.UserID, entity.PlaylistRoleViewer
	case *tg.ChannelParticipantCreator:
		return p.UserID, entity.PlaylistRoleOwner
	case *tg.ChannelParticipantAdmin:
		return p.UserID, entity.PlaylistRoleModerator
	default:
		return 0, ""
	}
}
