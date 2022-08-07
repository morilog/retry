package retry

import (
	"context"
	"time"
)

// Default Delay func options
const (
	DefaultAttempts    uint8 = 10
	DefaultDelay             = time.Millisecond * 100
	DefaultDelayFactor       = 1
)

type config struct {
	maxAttempts uint8
	delay       time.Duration
	delayFactor int
	stopRetryIf StopRetryIfFunc
	onRetry     OnRetryFunc
}

// Option is a function type for
// manipulating Retry configs
type Option func(*config)

// MaxAttempts is an option to set maximum attempts number
// It could be a integer number between 0-255
// Default is 10
func MaxAttempts(attempts uint8) Option {
	return func(c *config) {
		c.maxAttempts = attempts
	}
}

// Delay is an option to set minimum delay between each retries
// Default is 100ms
func Delay(delay time.Duration) Option {
	return func(c *config) {
		c.delay = delay
	}
}

// DelayFactor is an option to set increase of delay
// duration on each retries
// Actual delay calculated by: (delay * delayFactor * attemptNumber)
// Default is 1
func DelayFactor(factor int) Option {
	return func(c *config) {
		c.delayFactor = factor
	}
}

// StopRetryIfFunc is a function type to set conditions
// To stop continuing the retry mechanism
type StopRetryIfFunc func(ctx context.Context, err error) bool

// StopRetryIf is an option to set StopRetryIfFunc
func StopRetryIf(fn StopRetryIfFunc) Option {
	return func(c *config) {
		c.stopRetryIf = fn
	}
}

// OnRetryFunc is a function type to set some
// functionality to retry function
// It stops retry mechanism when error was not nil
type OnRetryFunc func(ctx context.Context, attempt int) error

// OnRetry is an option to set OnRetryFunc
func OnRetry(fn OnRetryFunc) Option {
	return func(c *config) {
		c.onRetry = fn
	}
}

// Retry tries to handle the operation func
// It tries until maxAttempts exceeded
// At the end of retries returns error from operation func
func Retry(ctx context.Context, operation func() error, opts ...Option) error {
	cfg := &config{
		maxAttempts: DefaultAttempts,
		delay:       DefaultDelay,
		delayFactor: DefaultDelayFactor,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	var err error
	for attempt := 0; attempt < int(cfg.maxAttempts); attempt++ {
		err = operation()
		if err == nil {
			return nil
		}

		if cfg.stopRetryIf != nil {
			if cfg.stopRetryIf(ctx, err) {
				break
			}
		}

		if cfg.onRetry != nil {
			if err := cfg.onRetry(ctx, attempt+1); err != nil {
				return err
			}
		}

		select {
		case <-ctx.Done():
			return err
		case <-time.After(cfg.delay * time.Duration((attempt+1)*cfg.delayFactor)):
			continue
		}
	}

	return err
}
