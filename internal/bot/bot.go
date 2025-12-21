package bot

import (
	appHandlers "backend/internal/bot/handlers"
	"backend/internal/infra"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/celestix/gotgproto/types"
	"github.com/gotd/td/telegram/dcs"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
)

func NewBotClient(cfg *infra.Config, logger *zap.Logger) (*gotgproto.Client, error) {
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
			Session:          sessionMaker.SqlSession(sqlite.Open("./telegram/bot.db")),
			Logger:           logger,
		},
	)
	if err != nil {
		return nil, err
	}

	self := client.Self
	logger.Info("bot logged in", zap.String("username", self.Username))

	go func() {
		// TODO: add stop with FX
		if err := client.Idle(); err != nil {
			logger.Error("failed to start a bot", zap.Error(err))
		}
	}()

	return client, nil
}

func RegisterHandlers(
	client *gotgproto.Client,
	groupHandler *appHandlers.GroupHandler,
	chatActionHandler *appHandlers.ChatActionHandler,
	startHandler *appHandlers.StartHandler,
) {
	client.Dispatcher.AddHandler(handlers.NewChatMemberUpdated(nil, groupHandler.HandleGroup))
	client.Dispatcher.AddHandler(handlers.NewCommand("start", startHandler.HandleStart))
	client.Dispatcher.AddHandler(handlers.NewMessage(func(msg *types.Message) bool {
		switch msg.Action.(type) {
		case *tg.MessageActionChatEditTitle,
			*tg.MessageActionChatAddUser,
			*tg.MessageActionChatDeleteUser,
			*tg.MessageActionChatJoinedByLink:
			return true
		default:
			return false
		}
	}, chatActionHandler.HandleChatAction))
}
