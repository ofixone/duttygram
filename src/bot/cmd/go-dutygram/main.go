package main

import (
	"context"
	"dutygram/cmd/go-dutygram/config"
	httpserver "dutygram/internal/api/http"
	logprettier "dutygram/pkg/log"
	tarantoolwrapper "dutygram/pkg/tarantool"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/jessevdk/go-flags"
	tarantooldriver "github.com/tarantool/go-tarantool"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type environment string

const (
	local environment = "local"
	dev   environment = "dev"
	prod  environment = "prod"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := initCfg()
	env := environment(cfg.Env)
	log := initLog(env)

	log.Debug("start init")

	var db *tarantoolwrapper.Wrapper
	var b *bot.Bot

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Debug("start telegram bot", slog.String(
			"listen",
			cfg.ListenHTTP,
		), slog.Any("telegram", cfg.Telegram))
		server := httpserver.NewServer(log, db)
		b = initTelegramBot(ctx, telegramOpts{
			token:      cfg.Telegram.Token,
			webhookUri: cfg.Telegram.WebhookUri,
		}, server.HandleWebhook, env == local)
		log.Debug("successfully started telegram bot")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Debug("start tarantool", slog.Any("tarantool", cfg.Tarantool))
		client := initTarantool(ctx, tarantoolOpts{
			listen:   cfg.Tarantool.Listen,
			username: cfg.Tarantool.Username,
			password: cfg.Tarantool.Password,
		})
		db = client
		log.Debug("successfully started tarantool")
	}()

	wg.Wait()

	log.Debug("init done, run webhook handler")
	go b.StartWebhook(ctx)
	err := http.ListenAndServe(cfg.ListenHTTP, b.WebhookHandler())
	if err != nil {
		notifyAndExitIfErr(err)
	}
}

func initCfg() *config.Config {
	cfg := &config.Config{}
	parser := flags.NewParser(cfg, flags.Default)
	parser.SubcommandsOptional = true
	_, parserErr := parser.AddCommand("version", "Show version", "Show build version", &struct{}{})
	notifyAndExitIfErr(parserErr)
	_, parserErr = parser.Parse()
	notifyAndExitIfErr(parserErr)

	return cfg
}

func initLog(env environment) *slog.Logger {
	var log *slog.Logger

	switch env {
	case local:
		opts := logprettier.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}

		handler := opts.NewPrettyHandler(os.Stdout)
		log = slog.New(handler)
	case dev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case prod:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func notifyAndExitIfErr(err error) {
	if err != nil {
		fmt.Printf("Error init: %s", err)
		os.Exit(1)
	}
}

type telegramOpts struct {
	token      string
	webhookUri string
}

func initTelegramBot(ctx context.Context, opts telegramOpts, f bot.HandlerFunc, debug bool) *bot.Bot {
	botOpts := []bot.Option{
		bot.WithDefaultHandler(f),
	}
	if debug {
		botOpts = append(botOpts, bot.WithDebug())
	}

	var err error
	b, err := bot.New(opts.token, botOpts...)
	if err != nil {
		panic(err)
	}

	h, err := b.GetWebhookInfo(ctx)
	if err != nil {
		panic(err)
	}
	if h.URL != opts.webhookUri {
		_, err = b.DeleteWebhook(ctx, &bot.DeleteWebhookParams{DropPendingUpdates: true})
		notifyAndExitIfErr(err)
		_, err = b.SetWebhook(ctx, &bot.SetWebhookParams{
			URL: opts.webhookUri,
		})
		notifyAndExitIfErr(err)
	}

	return b
}

type tarantoolOpts struct {
	listen   string
	username string
	password string
}

func initTarantool(
	ctx context.Context,
	opts tarantoolOpts,
) *tarantoolwrapper.Wrapper {
	drvOpts := tarantooldriver.Opts{
		Timeout:   5 * time.Second,
		Reconnect: 1 * time.Second,
		User:      opts.username,
		Pass:      opts.password,
	}

	client, err := tarantoolwrapper.ConnectWithRetries(ctx, opts.listen, drvOpts, tarantoolwrapper.ConnectRetryOpts{})
	if err != nil {
		notifyAndExitIfErr(fmt.Errorf("failed connect to tarantool: %w", err))
	}

	return tarantoolwrapper.New(client)
}
