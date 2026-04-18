// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coding-hui/iam/internal/api"
	"github.com/coding-hui/iam/internal/driver"
	"github.com/coding-hui/iam/pkg/shutdown"
	"github.com/coding-hui/iam/pkg/shutdown/shutdownmanagers/posixsignal"
)

// Run starts the API server.
func Run(basename string, opts *Options) error {
	cfg := &opts.Config

	// Create registry
	reg := driver.NewRegistry(cfg)
	logger := reg.Logger()

	// Initialize context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize all lazy services
	if err := reg.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}

	// Run database migrations
	if err := reg.MigrateUp(ctx); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Create Gin router
	router := api.NewRouter(reg)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Set up graceful shutdown
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	// Add shutdown callbacks
	gs.AddShutdownCallback(shutdown.ShutdownFunc(func(shutdownManager string) error {
		return srv.Shutdown(context.Background())
	}))

	gs.AddShutdownCallback(shutdown.ShutdownFunc(func(shutdownManager string) error {
		return reg.Persister().Close(ctx)
	}))

	gs.SetErrorHandler(shutdown.ErrorFunc(func(err error) {
		logger.Errorf("shutdown error: %v", err)
	}))

	// Start shutdown manager in background
	if err := gs.Start(); err != nil {
		return fmt.Errorf("failed to start shutdown manager: %w", err)
	}

	// Start HTTP server in goroutine
	go func() {
		logger.Infof("Starting API server %s on %s", basename, addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Trigger graceful shutdown
	gs.StartShutdown(posixsignal.NewPosixSignalManager())

	// Wait for callbacks to complete with timeout
	time.Sleep(5 * time.Second)

	logger.Info("Server stopped")

	return nil
}
