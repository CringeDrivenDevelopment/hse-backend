package spotify

import (
	"backend/internal/transport/api/dto"
	"context"

	"github.com/zmb3/spotify/v2"
)

func (s *API) Search(ctx context.Context, query string) ([]dto.Track, error) {
	resp, err := s.client.Search(ctx, query, spotify.SearchTypeTrack)
	if err != nil {
		return nil, err
	}

	result := make([]dto.Track, len(resp.Tracks.Tracks))
	for i, track := range resp.Tracks.Tracks {
		result[i] = dto.Track{
			Id:        track.ID.String(),
			Title:     track.Name,
			Authors:   getArtists(track.Artists),
			Thumbnail: track.Album.Images[0].URL,
			Explicit:  track.Explicit,
			Length:    int32(track.Duration) / 1000,
		}
	}

	return result, nil
}
