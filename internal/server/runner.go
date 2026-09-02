package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

type Runner struct {
	logger          *slog.Logger
	shutdownTimeout time.Duration
}

func NewRunner(logger *slog.Logger, shutdownTimeout time.Duration) *Runner {
	return &Runner{
		logger:          logger,
		shutdownTimeout: shutdownTimeout,
	}
}

func (r *Runner) Run(ctx context.Context, server *http.Server) error {
	startTime := time.Now()

	errCh := make(chan error, 1)
	go func() {
		r.logger.Info("server started", slog.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), r.shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			r.logger.Error("graceful shutdown failed", slog.String("error", err.Error()))
			return errors.Join(err, server.Close())
		}

		r.logger.Info("server stopped successfully", slog.Duration("run time", time.Since(startTime)))
		return nil
	}
}
