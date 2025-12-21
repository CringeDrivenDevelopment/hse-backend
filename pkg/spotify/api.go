package spotify

import (
	"backend/internal/infra"
	"context"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/fx"
	"golang.org/x/oauth2/clientcredentials"
)

type API struct {
	client *spotify.Client
	auth   *spotifyauth.Authenticator
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
		auth: spotifyauth.New(
			spotifyauth.WithRedirectURL(os.Getenv("SPOTIFY_REDIRECT_URL")),
			spotifyauth.WithScopes(
				spotifyauth.ScopePlaylistModifyPublic,
				spotifyauth.ScopePlaylistModifyPrivate,
				spotifyauth.ScopePlaylistReadPrivate,
			),
			spotifyauth.WithClientID(os.Getenv("SPOTIFY_CLIENT_ID")),
			spotifyauth.WithClientSecret(os.Getenv("SPOTIFY_CLIENT_SECRET")),
		),
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
