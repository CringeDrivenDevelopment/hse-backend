package spotify

import (
	"backend/internal/infra"
	"backend/internal/transport/api/dto"
	"context"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/fx"
	"golang.org/x/oauth2/clientcredentials"
)

type API struct {
	client *spotify.Client
}

func New(lc fx.Lifecycle, cfg *infra.Config) (*API, error) {
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	config := &clientcredentials.Config{
		ClientID:     cfg.SpotifyId,
		ClientSecret: cfg.SpotifySecret,
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctxWithCancel)
	if err != nil {
		cancel()
		return nil, err
	}

	httpClient := spotifyauth.New().Client(ctxWithCancel, token)
	result := &API{
		client: spotify.New(httpClient),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			cancel()
			return nil
		},
	})

	return result, nil
}

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
