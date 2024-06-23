package http

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
)

type Server struct {
	logger *slog.Logger
}

func NewServer(logger *slog.Logger) *Server {
	return &Server{logger: logger}
}

func (s *Server) HandleWebhook(_ context.Context, _ *bot.Bot, update *models.Update) {
	s.logger.Debug("http.Server.HandleWebhook", slog.Any("update", update))
}
