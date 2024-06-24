package config

type Config struct {
	Env        string `long:"env" description:"Environment" env:"ENV" default:"prod"`
	ListenHTTP string `long:"listen-http" description:"Listening host:port for http-server" env:"HTTP_SERVICE_LISTEN" required:"true"`
	Telegram   struct {
		Token      string `long:"telegram-token" description:"Telegram token" env:"TELEGRAM_TOKEN" required:"true"`
		WebhookUri string `long:"telegram-webhook-uri" description:"Telegram webhook uri" env:"TELEGRAM_WEBHOOK_URI" required:"true"`
	}
	Tarantool struct {
		Listen   string `long:"tarantool-listen" description:"Tarantool listen host:port" env:"TNT_LISTEN" required:"true"`
		Username string `long:"tarantool-user" description:"Tarantool username" env:"TNT_USERNAME" required:"true"`
		Password string `long:"tarantool-password" description:"Tarantool password" env:"TNT_PASSWORD" required:"true"`
	}
}
