package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ryoeuyo/ratelim/internal/config"
	"github.com/ryoeuyo/ratelim/internal/handler"
	"github.com/ryoeuyo/ratelim/internal/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load("config.yaml")
	if err != nil {
		slog.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.SlogLevel(),
	}))
	h := &handler.ProxyHandler{}
	srv := &http.Server{
		Handler:      h,
		Addr:         cfg.Server.Addr,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	runner := server.NewRunner(logger, cfg.Server.ShutdownTimeout)

	if err := runner.Run(ctx, srv); err != nil {
		logger.Error("error while running server", slog.String("error", err.Error()))
	}
}
