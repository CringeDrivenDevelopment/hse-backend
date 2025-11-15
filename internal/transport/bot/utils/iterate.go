package utils

import (
	"backend/internal/infra/queries"
	"backend/internal/transport/bot/models"
	"context"
	"errors"

	"github.com/celestix/gotgproto"
	"github.com/gotd/td/tg"
)

func iterateParticipants(client *gotgproto.Client, ctx context.Context, channel *tg.InputChannel) (*[]models.ParticipantData, error) {
	seen := make(map[int64]bool)
	var data []models.ParticipantData

	filters := []tg.ChannelParticipantsFilterClass{
		&tg.ChannelParticipantsRecent{},
		&tg.ChannelParticipantsAdmins{},
	}

	for _, filter := range filters {
		offset := 0
		const limit = 100

		for {
			resp, err := client.API().ChannelsGetParticipants(ctx, &tg.ChannelsGetParticipantsRequest{
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
				userId, role := extractUserAndRole(participant)
				if userId == 0 || userId == client.Self.ID || seen[userId] {
					continue
				}

				seen[userId] = true
				data = append(data, models.ParticipantData{
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

func extractUserAndRole(participant tg.ChannelParticipantClass) (int64, queries.PlaylistRole) {
	switch p := participant.(type) {
	case *tg.ChannelParticipant:
		return p.UserID, queries.PlaylistRoleViewer
	case *tg.ChannelParticipantCreator:
		return p.UserID, queries.PlaylistRoleOwner
	case *tg.ChannelParticipantAdmin:
		return p.UserID, queries.PlaylistRoleModerator
	default:
		return 0, ""
	}
}
