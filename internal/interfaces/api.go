package interfaces

import (
	"backend/internal/api/dto"
	"context"
)

// MusicAPI - common search API for multiple services
type MusicAPI interface {
	Search(ctx context.Context, query string) ([]dto.Track, error)
}
