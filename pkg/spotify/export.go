package spotify

import (
	"backend/internal/domain/entity"
	"backend/internal/transport/api/dto"
	"context"
	"errors"
	"time"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

func (s *API) Export(ctx context.Context, token *oauth2.Token, pl dto.Playlist, isPublic bool) error {
	// check playlist type
	if pl.Type != string(entity.PlaylistTypeSpotify) {
		return errors.New("invalid playlist type")
	}

	// Build description with timestamp indicating export source.
	desc := "Muse export " + time.Now().Format("2006-01-02 15:04:05")

	httpClient := s.auth.Client(context.Background(), token)
	spotifyClient := spotify.New(httpClient)

	// Get current user's ID
	user, err := spotifyClient.CurrentUser(context.Background())
	if err != nil {
		return err
	}

	// Create a new playlist for the authenticated user.
	// Public flag is set to true to allow optional public creation.
	userPlaylist, err := s.client.CreatePlaylistForUser(ctx, user.ID, pl.Title, desc, isPublic, false)
	if err != nil {
		return err
	}

	// Build a set of allowed IDs for quick lookup.
	allowed := make(map[string]struct{}, len(pl.AllowedIds))
	for _, id := range pl.AllowedIds {
		allowed[id] = struct{}{}
	}

	// Gather Spotify IDs that are both in the playlist tracks and allowed.
	var ids []spotify.ID
	for _, tr := range pl.Tracks {
		if _, ok := allowed[tr.Id]; ok {
			ids = append(ids, spotify.ID(tr.Id))
		}
	}

	// Spotify API accepts up to 100 tracks per request.
	const batchSize = 100
	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}
		batch := ids[i:end]
		if _, err = s.client.AddTracksToPlaylist(ctx, userPlaylist.ID, batch...); err != nil {
			return err
		}
	}

	return nil
}
