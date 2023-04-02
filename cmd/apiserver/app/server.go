package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/wecoding/iam/cmd/apiserver/app/options"
	"github.com/wecoding/iam/pkg/apiserver"
	"github.com/wecoding/iam/version"
)

// NewAPIServerCommand creates a *cobra.Command object with default parameters
func NewAPIServerCommand() *cobra.Command {
	s := options.NewServerRunOptions()
	cmd := &cobra.Command{
		Use: "apiserver",
		Long: `The IAM API apiserver validates and configures data for the API objects. 
The API Server services REST operations and provides the frontend to the
cluster's shared state through which all other components interact.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := s.Validate(); err != nil {
				return err
			}
			return Run(s)
		},
		SilenceUsage: true,
	}

	fs := cmd.Flags()
	namedFlagSets := s.Flags()
	for _, set := range namedFlagSets.FlagSets {
		fs.AddFlagSet(set)
	}

	return cmd
}

// Run runs the specified APIServer. This should never exit.
func Run(s *options.ServerRunOptions) error {
	// The apiserver is not terminal, there is no color default.
	// Force set to false, this is useful for the dry-run API.
	color.NoColor = false

	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := run(ctx, s, errChan); err != nil {
			errChan <- fmt.Errorf("failed to run apiserver: %w", err)
		}
	}()
	var term = make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		klog.Infof("Received SIGTERM, exiting gracefully...")
	case err := <-errChan:
		klog.Errorf("Received an error: %s, exiting gracefully...", err.Error())
		return err
	}
	klog.Infof("See you next time!")

	return nil
}

func run(ctx context.Context, s *options.ServerRunOptions, errChan chan error) error {
	klog.Infof("IAM information: version: %v, gitRevision: %v", version.IAMVersion, version.GitRevision)

	if s.GenericServerRunOptions.PrintVersion {
		version.PrintVersionAndExit()
	}

	server := apiserver.New(*s.GenericServerRunOptions)

	return server.Run(ctx, errChan)
}
