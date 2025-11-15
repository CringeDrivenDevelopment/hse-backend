package main

import (
	"backend/internal/infra"
	"backend/internal/service"
	"backend/internal/transport/api/handlers"
	"backend/internal/transport/api/middlewares"
	"backend/pkg/spotify"
	"backend/pkg/youtube"

	"go.uber.org/fx"
)

func main() {
	// TODO: log db requests
	// TODO: add otel
	// TODO: add image proxy, DL
	// TODO: process errors like .Error(), if code is 500 - print stacktrace

	fx.New(
		fx.Provide(
			// REST API
			infra.NewEcho,
			infra.NewHuma,
			middlewares.NewLogger,
			middlewares.NewAuth,
			handlers.NewAuth,
			handlers.NewPlaylist,
			handlers.NewTrack,

			// services and infra
			infra.NewLogger,
			infra.NewConfig,
			infra.NewPostgresConnection,
			youtube.New,
			spotify.New,
			service.NewAuth,
			service.NewPermission,
			service.NewPlaylist,
			service.NewTrack,
			service.NewUser,
		),
		fx.Invoke(func(auth *handlers.Auth, track *handlers.Track, playlist *handlers.Playlist) {
			// need each of controllers, to register them

			// no need to call infra, apis and services, they're deps, started automatically
		}),
	).Run()
}
