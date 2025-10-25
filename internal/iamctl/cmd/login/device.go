// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package login

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/coding-hui/wecoding-sdk-go/rest"

	"github.com/coding-hui/iam/internal/iamctl/cmd/util"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/cli/genericclioptions"
)

// DeviceLoginOptions defines the options for device login.
type DeviceLoginOptions struct {
	ClientID  string
	Scope     string
	ServerURL string
	IOStreams genericclioptions.IOStreams

	// Internal state
	restClient *rest.RESTClient
}

// NewDeviceLoginOptions creates a new DeviceLoginOptions instance.
func NewDeviceLoginOptions(ioStreams genericclioptions.IOStreams) *DeviceLoginOptions {
	return &DeviceLoginOptions{
		IOStreams: ioStreams,
	}
}

// NewCmdDeviceLogin creates a new device login command.
func NewCmdDeviceLogin(f util.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewDeviceLoginOptions(ioStreams)

	cmd := &cobra.Command{
		Use:   "device",
		Short: "Login using OAuth 2.0 Device Authorization Grant",
		Long:  "Login using OAuth 2.0 Device Authorization Grant flow",
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(o.Complete(f, cmd, args))
			util.CheckErr(o.Validate())
			util.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVar(&o.ClientID, "client-id", "iamctl", "OAuth client ID")
	cmd.Flags().StringVar(&o.Scope, "scope", "", "OAuth scope")
	cmd.Flags().StringVar(&o.ServerURL, "server", "http://localhost:8000", "IAM server URL")

	return cmd
}

// Complete completes the required options.
func (o *DeviceLoginOptions) Complete(f util.Factory, cmd *cobra.Command, args []string) error {
	config, err := f.ToRESTConfig()
	if err != nil {
		return fmt.Errorf("failed to get REST config: %w", err)
	}

	// Override server URL if provided
	if o.ServerURL != "" {
		config.Host = o.ServerURL
	}

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		return fmt.Errorf("failed to create REST client: %w", err)
	}

	o.restClient = restClient
	return nil
}

// Validate validates the provided options.
func (o *DeviceLoginOptions) Validate() error {
	if o.ClientID == "" {
		return fmt.Errorf("client-id is required")
	}
	if o.ServerURL == "" {
		return fmt.Errorf("server URL is required")
	}
	return nil
}

// Run executes the device login command.
func (o *DeviceLoginOptions) Run() error {
	fmt.Fprintf(o.IOStreams.Out, "Starting OAuth 2.0 Device Authorization Grant flow...\n")
	fmt.Fprintf(o.IOStreams.Out, "Client ID: %s\n", o.ClientID)
	if o.Scope != "" {
		fmt.Fprintf(o.IOStreams.Out, "Scope: %s\n", o.Scope)
	}

	// Step 1: Request device authorization
	fmt.Fprintf(o.IOStreams.Out, "\n1. Requesting device authorization...\n")

	authReq := v1.DeviceAuthorizationRequest{
		ClientID: o.ClientID,
		Scope:    o.Scope,
	}

	var authResp v1.DeviceAuthorizationResponse
	err := o.restClient.Post().
		AbsPath("/api/v1/device/code").
		Body(&authReq).
		Do(context.TODO()).
		Into(&authResp)

	if err != nil {
		return fmt.Errorf("failed to request device authorization: %w", err)
	}

	fmt.Fprintf(o.IOStreams.Out, "✓ Device authorization requested successfully\n")
	fmt.Fprintf(o.IOStreams.Out, "\nPlease visit this URL in your browser:\n")
	fmt.Fprintf(o.IOStreams.Out, "%s\n", authResp.VerificationURIComplete)
	fmt.Fprintf(o.IOStreams.Out, "\nAnd enter the following user code:\n")
	fmt.Fprintf(o.IOStreams.Out, "%s\n", authResp.UserCode)
	fmt.Fprintf(o.IOStreams.Out, "\nWaiting for user authorization...\n")

	// Step 2: Poll for token
	fmt.Fprintf(o.IOStreams.Out, "\n2. Polling for access token...\n")

	for i := 0; i < 120; i++ { // Max 10 minutes (600 seconds / 5 seconds interval)
		time.Sleep(time.Duration(authResp.Interval) * time.Second)

		tokenReq := v1.DeviceTokenRequest{
			DeviceCode: authResp.DeviceCode,
			ClientID:   o.ClientID,
		}

		var tokenResp v1.DeviceTokenResponse
		err = o.restClient.Post().
			AbsPath("/api/v1/device/token").
			Body(&tokenReq).
			Do(context.TODO()).
			Into(&tokenResp)

		if err == nil {
			fmt.Fprintf(o.IOStreams.Out, "✓ Authorization granted!\n")
			fmt.Fprintf(o.IOStreams.Out, "\nAccess Token: %s\n", tokenResp.AccessToken)
			fmt.Fprintf(o.IOStreams.Out, "Token Type: %s\n", tokenResp.TokenType)
			fmt.Fprintf(o.IOStreams.Out, "Expires In: %d seconds\n", tokenResp.ExpiresIn)
			if tokenResp.Scope != "" {
				fmt.Fprintf(o.IOStreams.Out, "Scope: %s\n", tokenResp.Scope)
			}

			fmt.Fprintf(o.IOStreams.Out, "\nYou can now use this access token to authenticate with IAM.\n")
			return nil
		}

		// Check if authorization is still pending
		if i%6 == 0 { // Print status every 30 seconds
			fmt.Fprintf(o.IOStreams.Out, ".")
		}
	}

	return fmt.Errorf("authorization timeout - please try again")
}
