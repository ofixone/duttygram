package main

import (
	"context"
	httpserver "dutygram/api/http"
	"dutygram/cmd/go-dutygram/config"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/jessevdk/go-flags"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

var b *bot.Bot

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := &config.Config{}
	parser := flags.NewParser(cfg, flags.Default)
	parser.SubcommandsOptional = true
	_, parserErr := parser.AddCommand("version", "Show version", "Show build version", &struct{}{})
	notifyAndExitIfErr(parserErr)
	_, parserErr = parser.Parse()
	notifyAndExitIfErr(parserErr)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		startTelegramBot(ctx, cfg.ListenHTTP, cfg.Telegram.Token, cfg.Telegram.WebhookUri)
	}()

	wg.Wait()
	b.StartWebhook(ctx)
}

func notifyAndExitIfErr(err error) {
	if err != nil {
		fmt.Printf("Error init: %s.\nFor help use -h\n", err)
		os.Exit(1)
	}
}

func startTelegramBot(ctx context.Context, listen string, token string, webhookUri string) {
	server := httpserver.NewServer(slog.Default())
	opts := []bot.Option{
		bot.WithDefaultHandler(server.HandleWebhook),
		bot.WithDebug(),
	}

	var err error
	b, err = bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	h, err := b.GetWebhookInfo(ctx)
	if err != nil {
		panic(err)
	}
	if h.URL != webhookUri {
		_, err = b.DeleteWebhook(ctx, &bot.DeleteWebhookParams{DropPendingUpdates: true})
		notifyAndExitIfErr(err)
		_, err = b.SetWebhook(ctx, &bot.SetWebhookParams{
			URL: webhookUri,
		})
		notifyAndExitIfErr(err)
	}

	go func() {
		err := http.ListenAndServe(listen, b.WebhookHandler())
		if err != nil {
			return
		}
	}()
}
