package main

import (
	"backend/internal/bot"
	"backend/internal/bot/handlers"
	"backend/internal/bot/services"
	"backend/internal/infra"
	"backend/internal/service"
	"log"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/gotd/td/telegram/dcs"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
)

func main() {
	cfg, err := infra.NewConfig()
	if err != nil {
		log.Fatal("failed to load config", err)
	}

	logger, err := infra.NewLogger(cfg.Debug)
	if err != nil {
		log.Fatal("failed to create logger", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Println(err)
		}
	}()

	db, err := infra.NewPostgres(cfg.DbUrl)
	if err != nil {
		logger.Fatal("failed to connect to postgres", zap.Error(err))
	}

	playlistService := service.NewPlaylistService(db, logger)
	permissionService := service.NewPermissionService(db, logger)

	var dcList dcs.List
	if cfg.Debug {
		dcList = dcs.Test()
	} else {
		dcList = dcs.Prod()
	}

	client, err := gotgproto.NewClient(
		cfg.AppId,
		cfg.AppHash,
		gotgproto.ClientTypeBot(cfg.BotToken),
		&gotgproto.ClientOpts{
			DCList:           dcList,
			DisableCopyright: true,
			InMemory:         true,
			Session:          sessionMaker.SqlSession(sqlite.Open("telegram/bot.db")),
			Logger:           logger,
		},
	)
	if err != nil {
		logger.Fatal("failed to create a bot client", zap.Error(err))
	}

	self, err := client.GetSelf()
	if err != nil {
		logger.Fatal("failed to get self", zap.Error(err))
	}
	logger.Info("bot logged in", zap.String("username", self.Username))

	chatService := services.NewChatService(client)
	partService := services.NewParticipantService()

	groupHandler := handlers.NewGroupHandler(playlistService, permissionService, chatService, partService, client, logger)
	chatActionHandler := handlers.NewChatActionHandler(playlistService, logger)
	startHandler := handlers.NewStartHandler(logger)

	bot.RegisterHandlers(client, groupHandler, chatActionHandler, startHandler)

	if err = client.Idle(); err != nil {
		logger.Fatal("failed to start a bot", zap.Error(err))
	}
}
