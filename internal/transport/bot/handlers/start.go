package handlers

import (
	"github.com/celestix/gotgproto/ext"
)

func (b *Bot) handleStart(ctx *ext.Context, update *ext.Update) error {
	_, err := ctx.Reply(update, ext.ReplyTextString("Привет, я Лотти - бот для модерации плейлистов."+"\n\n"+
		"Добавь меня в группу и я подгружу данные из неё"+"\n\n"+
		"Для управления плейлистами, к которым у тебя есть доступ - зайди в миниапп)"), &ext.ReplyOpts{})
	if err != nil {
		b.logger.Error(err.Error())
		return nil
	}

	return nil
}
