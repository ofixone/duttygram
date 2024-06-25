package telegram

import (
	"context"
	"github.com/go-telegram/bot"
)

type Messenger struct {
	api *bot.Bot
}

func NewMessenger(api *bot.Bot) *Messenger {
	return &Messenger{api: api}
}

func (m Messenger) SendText(ctx context.Context, content string) error {
	return nil
}
