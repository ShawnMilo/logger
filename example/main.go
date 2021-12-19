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
	fmt.Printf("user_id: %s\n", userID)
	crash(ctx)
}

func crash(ctx context.Context) {
	lg := logger.FromContext(ctx)
	lg.Error("broken")
}
