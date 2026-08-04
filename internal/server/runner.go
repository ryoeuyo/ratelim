package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Runner struct {
	logger *slog.Logger
}

func NewRunner(logger *slog.Logger) Runner {
	return Runner{
		logger: logger,
	}
}

func (r *Runner) Run(ctx context.Context, server *http.Server) error {
	startTime := time.Now()

	r.logger.Info("server started", slog.String("address", server.Addr))
	errCh := make(chan error)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		r.logger.Info("server stopped successfully", slog.Duration("run time", time.Since(startTime)))
		return nil
	}
}
