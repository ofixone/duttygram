package http

import (
	"context"
	tarantoolwrapper "dutygram/pkg/tarantool"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
)

type Server struct {
	logger *slog.Logger
	db     *tarantoolwrapper.Wrapper
}

func NewServer(logger *slog.Logger, db *tarantoolwrapper.Wrapper) *Server {
	return &Server{logger: logger, db: db}
}

func (s *Server) HandleWebhook(_ context.Context, _ *bot.Bot, update *models.Update) {
	s.logger.Debug("http.Server.HandleWebhook", slog.Any("update", update))
}
