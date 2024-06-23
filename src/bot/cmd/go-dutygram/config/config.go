package config

type Config struct {
	ListenHTTP string `long:"listen-http" description:"Listening host:port for http-server" env:"CL_HTTP_SERVICE_LISTEN" required:"true"`

	Telegram struct {
		Token      string `long:"telegram-token" description:"Telegram token" env:"CL_TELEGRAM_TOKEN" required:"true"`
		WebhookUri string `long:"telegram-webhook-uri" description:"Telegram webhook uri" env:"CL_TELEGRAM_WEBHOOK_URI" required:"true"`
	}
}
