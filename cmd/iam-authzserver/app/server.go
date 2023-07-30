// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/iam/cmd/iam-authzserver/app/options"
	"github.com/coding-hui/iam/internal/authzserver"
	"github.com/coding-hui/iam/internal/authzserver/config"
	"github.com/coding-hui/iam/pkg/app"
)

const commandDesc = `Authorization server to run casbin policies which can protecting your resources.

It provides support for enforcing authorization based on various access control models.

Find more casbin information at:
    https://github.com/casbin/casbin`

// NewAuthzServerAPP creates a *app.App object with default parameters.
func NewAuthzServerAPP(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("IAM Authz Server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(Run(opts)),
	)

	return application
}

// Run runs the specified AuthzServer. This should never exit.
func Run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		errChan := make(chan error)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			if err := run(ctx, opts, errChan); err != nil {
				errChan <- fmt.Errorf("failed to run authzserver: %w", err)
			}
		}()
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)

		select {
		case <-term:
			log.Infof("Received SIGTERM, exiting gracefully...")
		case err := <-errChan:
			log.Errorf("Received an error: %s, exiting gracefully...", err.Error())
			return err
		}
		log.Infof("See you next time!")

		return nil
	}
}

func run(ctx context.Context, opts *options.Options, errChan chan error) error {
	cfg, err := config.CreateConfigFromOptions(opts)
	if err != nil {
		return err
	}

	server, err := authzserver.New(cfg)
	if err != nil {
		return err
	}

	return server.Run(ctx, errChan)
}
