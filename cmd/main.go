package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/ryoeuyo/ratelim/internal/handler"
	"github.com/ryoeuyo/ratelim/internal/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	h := &handler.ProxyHandler{}
	srv := &http.Server{
		Handler: h,
		Addr:    "localhost:8080",
	}

	runner := server.NewRunner(logger)

	if err := runner.Run(ctx, srv); err != nil {
		logger.Error("error while running server", slog.String("error", err.Error()))
	}
}
