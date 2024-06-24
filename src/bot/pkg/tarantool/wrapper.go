package tarantool

import (
	"context"
	"fmt"
	"github.com/tarantool/go-tarantool"
	"time"
)

type Client interface {
	Close(ctx context.Context) error
	Ping(ctx context.Context) (*tarantool.Response, error)
	Call17(ctx context.Context, functionName string, args interface{}) (*tarantool.Response, error)
	Call17Typed(ctx context.Context, functionName string, args interface{}, result interface{}) error
	Call17Async(ctx context.Context, functionName string, args interface{}) *tarantool.Future
}

type Wrapper struct {
	client tarantool.Connector
}

func New(client tarantool.Connector) *Wrapper {
	return &Wrapper{client: client}
}

func (c *Wrapper) Close(ctx context.Context) error {
	var err error
	requestDoneCh := make(chan struct{})
	go func() {
		defer close(requestDoneCh)
		err = c.client.Close()
	}()

	select {
	case <-requestDoneCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Wrapper) Ping(ctx context.Context) (*tarantool.Response, error) {
	requestDoneCh := make(chan struct{})
	var resp *tarantool.Response
	var err error
	go func() {
		defer close(requestDoneCh)
		resp, err = c.client.Ping()
	}()

	select {
	case <-requestDoneCh:
		return resp, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Wrapper) Call17(
	ctx context.Context,
	functionName string,
	args interface{},
) (*tarantool.Response, error) {
	future := c.client.Call17Async(functionName, args)

	select {
	case <-future.WaitChan():
		return future.Get()
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Wrapper) Call17Typed(
	ctx context.Context,
	functionName string,
	args interface{},
	result interface{},
) error {
	future := c.client.Call17Async(functionName, args)

	select {
	case <-future.WaitChan():
		return future.GetTyped(result)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Wrapper) Call17Async(_ context.Context, functionName string, args interface{}) *tarantool.Future {
	return c.client.Call17Async(functionName, args)
}

// ConnectRetryOpts опции для подключения в несколько попыток
type ConnectRetryOpts struct {
	// Количество попыток подключения, по-умолчанию 10
	Attempts int
	// Пауза между попытками подключения, по-умолчанию 1 секунда
	Delay time.Duration
}

// ConnectWithRetries подключается к тарантулу в несколько попыток.
//
// Данный метод нужно использовать, когда ваше приложение не может сразу подключиться к тарантулу после запуска (не
// готова сеть, не подняты какие-нибудь sidecar контейнеры и т.п.)
func ConnectWithRetries(
	ctx context.Context,
	addr string,
	opts tarantool.Opts,
	retryOpts ConnectRetryOpts,
) (*tarantool.Connection, error) {
	if retryOpts.Attempts == 0 {
		retryOpts.Attempts = 10
	}
	if retryOpts.Delay == 0 {
		retryOpts.Delay = 1 * time.Second
	}

	var lastErr error
	for i := 0; i < retryOpts.Attempts; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		client, err := tarantool.Connect(addr, opts)
		lastErr = err
		if err != nil {
			select {
			case <-ctx.Done():
			case <-time.After(retryOpts.Delay):
			}
			continue
		}

		return client, nil
	}

	return nil, fmt.Errorf("can't establish connection to tarantool for %d attempts: %w",
		retryOpts.Attempts, lastErr)
}
