# logger

A simple JSON logger for Go.

It uses a `context.Context` to store values which will then be logged along with each message.

It is possible to recover these values but this should not be used to pass arguments into functions. Do not use this to pass required arguments.

## Usage

The API is very small.

1. Create a logger with `logger.New()`
1. (Optional) Call `SetDebug(bool)` to enable debug logging, or use the `DEBUG` environment variable
1. Log using `Debug`, `Info`, `Error`, `Debugf`, `Infof`, or `Errorf`
1. (Optional) As desired, use `With` to add a value to the loggers "tags"
1. (Optional) Recover a value from the tags using `ValueString`

```go
package main

import (
	"context"
	"fmt"

	"github.com/ShawnMilo/logger"
)

func main() {
	lg := logger.New()
	lg.Info("first message")
	ctx := context.Background()
	ctx = lg.With(ctx, "user_id", "123")
	lg.Info("second message")
	stuff(ctx)
}

func stuff(ctx context.Context) {
	lg := logger.FromContext(ctx)
	ctx = lg.With(ctx, "function", "stuff")
	lg.Info("thing message")
	userID := lg.ValueString("user_id")
}
```

**Output**:

```
{"level":"INFO","event_time":"2021-12-19T03:16:16Z","message":"first message"}
{"level":"INFO","event_time":"2021-12-19T03:16:16Z","message":"second message","tags":{"user_id":"123"}}
{"level":"INFO","event_time":"2021-12-19T03:16:16Z","message":"thing message","tags":{"function":"stuff","user_id":"123"}}
user_id: 123
```
