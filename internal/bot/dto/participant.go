package dto

import (
	"backend/internal/domain/entity"

	"github.com/gotd/td/tg"
)

type Participant struct {
	PrevRole entity.PlaylistRole
	NewRole  entity.PlaylistRole
	UserID   int64
	ChatID   int64
	ActorID  int64
}

type ChatInfo struct {
	Title string
	Photo tg.ChatPhotoClass
	Users *[]Participant
}
