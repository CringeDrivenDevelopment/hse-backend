package bot

import (
	appHandlers "backend/internal/bot/handlers"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/types"
	"github.com/gotd/td/tg"
)

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
