package interfaces

import (
	"backend/internal/transport/api/dto"
	"context"

	"golang.org/x/oauth2"
)

// SearchAPI - common search API for multiple services
type SearchAPI interface {
	Search(ctx context.Context, query string) ([]dto.Track, error)
}

// SpotifyAPI - spotify specific API for exporting data
type SpotifyExportAPI interface {
	Export(ctx context.Context, token *oauth2.Token, pl dto.Playlist, isPublic bool) error
}
