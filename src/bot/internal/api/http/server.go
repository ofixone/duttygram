package http

import (
	"context"
	"duttygram/internal"
	"duttygram/internal/telegram"
	tarantoolwrapper "duttygram/pkg/tarantool"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
	"strconv"
)

type Server struct {
	logger       *slog.Logger
	db           *tarantoolwrapper.Wrapper
	stateFactory *internal.StateFactory
}

func NewServer(logger *slog.Logger, stateFactory *internal.StateFactory) *Server {
	return &Server{logger: logger, stateFactory: stateFactory}
}

func (s *Server) HandleTelegramWebhook(ctx context.Context, bot *bot.Bot, update *models.Update) {
	state, err := s.stateFactory.Create(ctx, strconv.Itoa(int(update.Message.Chat.ID)), telegram.NewMessenger(bot))
	_ = state
	if err != nil {
		s.logger.Error("can't create state", err)
		return
	}
}
