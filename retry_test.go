package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/morilog/retry"
	"github.com/stretchr/testify/require"
)

var errProcess = errors.New("process error")
var errAnother = errors.New("another error")

func TestRetry(t *testing.T) {
	ctx := context.Background()
	retries := 0
	process := func() error {
		if retries <= 3 {
			retries++
			return errProcess
		}

		return nil
	}

	t.Run("Should failed and return error", func(t *testing.T) {
		attempts := 0
		err := retry.Retry(ctx, process, retry.MaxAttempts(1), retry.OnRetry(func(ctx context.Context, attempt int) error {
			attempts = attempt
			return nil
		}))
		require.NotNil(t, err)

		t.Run("Should attempts once", func(t *testing.T) {
			require.Equal(t, 1, attempts)
		})
	})

	t.Run("Should succeed on third try", func(t *testing.T) {
		attempts := 0
		err := retry.Retry(ctx, process, retry.MaxAttempts(4), retry.OnRetry(func(ctx context.Context, attempt int) error {
			attempts = attempt
			return nil
		}))
		require.Nil(t, err)

		t.Run("Should attempts three three times", func(t *testing.T) {
			require.Equal(t, 3, attempts)
		})
	})

	t.Run("Should attempts 10 times as default", func(t *testing.T) {
		attempts := 0
		err := retry.Retry(ctx, func() error {
			return errProcess
		}, retry.Delay(10*time.Millisecond), retry.OnRetry(func(ctx context.Context, attempt int) error {
			attempts = attempt
			return nil
		}))

		require.NotNil(t, err)
		require.Equal(t, 10, attempts)
	})

	t.Run("Should attempts for 450 millisecond", func(t *testing.T) {
		start := time.Now()
		_ = retry.Retry(ctx, func() error {
			return errProcess
		}, retry.Delay(time.Millisecond*10))
		howLong := time.Since(start)

		require.True(t, howLong >= 450*time.Millisecond)
		require.True(t, howLong <= 500*time.Millisecond)
	})

	t.Run("Should stop when error is not process error", func(t *testing.T) {
		attempts := 0
		err := retry.Retry(ctx, func() error {
			if attempts < 2 {
				attempts++
				return errProcess
			}

			return errAnother
		}, retry.StopRetryIf(func(ctx context.Context, err error) bool {
			if errors.Is(err, errAnother) {
				return true
			}

			return false
		}))

		require.NotNil(t, err)
		require.ErrorIs(t, err, errAnother)
		require.Equal(t, attempts, 2)
	})
}
