package interfaces

import (
	apiDto "backend/internal/api/dto"
	botDto "backend/internal/bot/dto"
	"backend/internal/domain/entity"
	"context"

	"golang.org/x/oauth2"
)

type UserService interface {
	Create(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) error
}

type AuthService interface {
	VerifyToken(authHeader string) (int64, error)
	GenerateToken(userID int64) (string, error)
	ParseTelegramData(initDataRaw string) (int64, error)
	ParseSpotifyData(spotifyAuthHeader string) (*oauth2.Token, error)
}

type PlaylistService interface {
	Create(ctx context.Context, title string, playlistType entity.PlaylistType, telegramId int64) (apiDto.Playlist, error)
	GetByGroup(ctx context.Context, telegramId int64) (apiDto.Playlist, error)
	GetById(ctx context.Context, playlistId string, userId int64) (apiDto.Playlist, error)
	GetAll(ctx context.Context, userId int64) ([]apiDto.Playlist, error)
	Rename(ctx context.Context, playlistId string, title string, userId int64) error
	UpdatePhoto(ctx context.Context, playlistId string, thumbnail string, userId int64) error
	Delete(ctx context.Context, playlistId string) error
}

type PermissionService interface {
	Add(ctx context.Context, role entity.PlaylistRole, playlist string, userId int64) error
	AddGroup(ctx context.Context, playlist string, users []botDto.Participant) error
	Remove(ctx context.Context, playlist string, userId int64) error
	Edit(ctx context.Context, role entity.PlaylistRole, playlist string, userId int64) error
	Get(ctx context.Context, userId int64, role entity.PlaylistRole) (string, error)
}

type TrackService interface {
	Search(ctx context.Context, platform, query string) ([]apiDto.Track, error)
	GetById(ctx context.Context, id string) (apiDto.Track, error)
	Approve(ctx context.Context, playlistId string, trackId string, userId int64) error
	Decline(ctx context.Context, playlistId string, trackId string, userId int64) error
	Submit(ctx context.Context, playlistId string, trackId string, userId int64) error
	Unapprove(ctx context.Context, playlistId string, trackId string, userId int64) error
}

type ParticipantService interface {
	Add(ctx context.Context, chatID int64, userIDs []int64) error
	Remove(ctx context.Context, chatID int64, userID int64) error
}
