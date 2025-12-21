package main

import (
	"backend/internal/bot"
	botHandlers "backend/internal/bot/handlers"
	botService "backend/internal/bot/service"
	"backend/internal/infra"
	"backend/internal/infra/repo"
	"backend/internal/service"
	"backend/internal/transport/api/handlers"
	"backend/internal/transport/api/middlewares"
	"backend/pkg/spotify"
	"backend/pkg/youtube"
	"fmt"
	"os"
	"time"

	"github.com/celestix/gotgproto"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	// TODO: log db requests
	// TODO: add otel
	// TODO: add image proxy, DL
	// TODO: process errors like .Error(), if code is 500 - print stacktrace
	cfg, err := infra.NewConfig()
	if err != nil {
		println(fmt.Sprintf("{\"level\":\"error\",\"ts\":%d,\"caller\":\"cmd/main.go\",\"msg\":\"%s\"}", time.Now().Unix(), err.Error()))
		os.Exit(1)
	}

	logger, err := infra.NewLogger(cfg)
	if err != nil {
		println(fmt.Sprintf("{\"level\":\"error\",\"ts\":%d,\"caller\":\"cmd/main.go\",\"msg\":\"%s\"}", time.Now().Unix(), err.Error()))
		os.Exit(1)
	}

	fx.New(
		fx.Supply(logger.Zap, logger, cfg),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			defer func(Zap *zap.Logger) {
				err := Zap.Sync()
				if err != nil {
					println(fmt.Sprintf("{\"level\":\"error\",\"ts\":%d,\"caller\":\"cmd/main.go\",\"msg\":\"%s\"}", time.Now().Unix(), err.Error()))
				}
			}(logger.Zap)

			return &fxevent.ZapLogger{Logger: logger.Zap}
		}),
		fx.Provide(
			// REST API
			middlewares.NewLogger,
			middlewares.NewAuth,
			handlers.NewAuth,
			handlers.NewPlaylist,
			handlers.NewTrack,

			// bot
			bot.NewBotClient,

			// bot handlers
			botHandlers.NewStartHandler,
			botHandlers.NewChatActionHandler,
			botHandlers.NewGroupHandler,

			// bot services
			botService.NewParticipantService,
			botService.NewChatService,

			// infra
			infra.NewEcho,
			infra.NewHuma,
			infra.NewPostgresConnection,

			// repositories
			repo.NewPlaylistRepo,
			repo.NewTrackRepo,
			repo.NewUserRepo,
			repo.NewPermissionRepo,

			// external APIs
			youtube.New,
			spotify.New,

			// service layer
			service.NewAuthService,
			service.NewPermissionService,
			service.NewPlaylistService,
			service.NewTrackService,
			service.NewUserService,
			service.NewParticipantService,
		),
		// bot.Module,
		fx.Invoke(func(
			auth *handlers.Auth,
			track *handlers.Track,
			playlist *handlers.Playlist,
			client *gotgproto.Client,
			groupHandler *botHandlers.GroupHandler,
			chatActionHandler *botHandlers.ChatActionHandler,
			startHandler *botHandlers.StartHandler) {
			// need each of controllers, to register them

			// no need to call infra, apis and services, they're deps, started automatically

			// register bot handlers
			// TODO: move to modules
			bot.RegisterHandlers(client, groupHandler, chatActionHandler, startHandler)
		}),
	).Run()
}
