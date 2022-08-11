![Tests](https://github.com/morilog/retry/actions/workflows/go-test.yml/badge.svg)

# Retry
Retry is a golang pure package implements the `retry-pattern` to retry failures until certain attempts in certain time.

## Installation
At first you need to download the library using `go get`
```bash
$ go get -u github.com/morilog/go-retry
```

Or just use it on your package by importing this:
```golang
import "github.com/morilog/retry"
```
And download it with `go mod`
```bash
$ go mod tidy
```

## Quick Start
Assume you have a function that calls an external api:

```golang
func getSomeExternalData(url string) error {
    // your business logic goes here
}
```
It might be failed during network errors or any other external errors. With `retry` you can retry your logic until maximum number of attempts exceeded.
```golang
package main

import "github.com/morilog/retry"

func main() {
    ctx := context.Background()

    // As default
    // It retries 10 times with increased 100ms delay
    err := retry.Retry(ctx, func() error {
        return getSomeExternalData("https://example.com")
    }, retry.Delay(time.Second))

    if err != nil {
        log.Fatal(err)
    }
}
```

## Options
### retry.Delay(delay time.Duration)
To set minimum delay duration between each retry.
> Default: 100ms
### retry.DelayFactor(factor int)
To set number that multiplies to delay to increase delay time by increasing attempts.
The delay for each retry calculated by this formula:
```
    AttemptDelay = Attempt * DelayFactor * Delay
```
> Default: 1
### retry.MaxAttempts(attempts uint8)
Sets maximum number of attempts to handle the operation

> Default: 10
### retry.StopRetryIf(fn StopRetryIfFunc)
Sometimes you need to cancel the retry mechanism when your operation returns some specific error. To achieve to this purpose, You can use this option.
> Default: `<not set>`

For example i stop the retry when an `SomeError` occurred:
```golang
err := retry.Retry(ctx, func() error {
    return getSomeExternalData("https://example.com")
},
    retry.Delay(time.Second),
    retry.StopRetryIf(func(ctx context.Context, err error) bool {
        if _, ok := err.(*SomeError); ok {
            return true
        }

        return false
    }),
)

if err != nil {
    log.Fatal(err)
}
```
### retry.OnRetry(fn OnRetryFunc)
Sometimes you need to know retry is on which attempt to do something like logging
> Default: `<not set>`

For example:
```golang
err := retry.Retry(ctx, func() error {
    return getSomeExternalData("https://example.com")
},
    retry.Delay(time.Second),
    retry.StopRetryIf(func(ctx context.Context, attempt int) error {
        logrus.Printf("The function retries for %d times.\n", attempt)
        return nil
    }),
)

if err != nil {
    log.Fatal(err)
}
```