package handlers

import (
	"github.com/celestix/gotgproto/ext"
	"go.uber.org/zap"
)

type StartHandler struct {
	logger *zap.Logger
}

func NewStartHandler(logger *zap.Logger) *StartHandler {
	return &StartHandler{logger: logger}
}

func (h *StartHandler) HandleStart(ctx *ext.Context, update *ext.Update) error {
	_, err := ctx.Reply(update, ext.ReplyTextString("Привет, я Лотти - бот для модерации плейлистов."+"\n\n"+
		"Добавь меня в группу и я подгружу данные из неё"+"\n\n"+
		"Для управления плейлистами, к которым у тебя есть доступ - зайди в миниапп)"), &ext.ReplyOpts{})
	if err != nil {
		h.logger.Error("failed to send start message", zap.Error(err))
		return nil
	}

	return nil
}
