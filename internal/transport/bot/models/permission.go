package models

import (
	"backend/internal/domain/entity"

	"github.com/gotd/td/tg"
)

type ParticipantData struct {
	PrevRole entity.PlaylistRole
	NewRole  entity.PlaylistRole
	UserID   int64
	ChatID   int64
	ActorID  int64
}

type Chat struct {
	Title string
	Photo tg.ChatPhotoClass
	Users *[]ParticipantData
}
