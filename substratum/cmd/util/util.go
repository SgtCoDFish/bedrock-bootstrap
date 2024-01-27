package util

import (
	"context"
	"log/slog"
	"os"
)

type key int

const (
	loggerKey key = 0
)

func Context() context.Context {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	return WithLogger(ctx, logger)
}

func Logger(ctx context.Context) *slog.Logger {
	return ctx.Value(loggerKey).(*slog.Logger)
}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
