// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	v1 "github.com/coding-hui/iam/pkg/api/authzserver/v1"
	"github.com/coding-hui/iam/pkg/api"
	cmdutil "github.com/coding-hui/iam/internal/iamctl/cmd/util"
	"github.com/coding-hui/iam/internal/iamctl/util/templates"
	"github.com/coding-hui/iam/pkg/cli/genericclioptions"
)

// CheckOptions is an options struct to support authz check subcommand.
type CheckOptions struct {
	Subject  string
	Resource string
	Action   string

	IOStreams genericclioptions.IOStreams
}

var checkExample = templates.Examples(`
	# Check if user-123 can perform GET on /api/users
	iamctl authz check user-123 /api/users GET

	# Check if user-456 can perform POST on /api/users:admin
	iamctl authz check user-456 /api/users:admin POST`)

// NewCmdCheck returns new initialized instance of check sub command.
func NewCmdCheck(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := &CheckOptions{IOStreams: ioStreams}

	cmd := &cobra.Command{
		Use:     "check SUBJECT RESOURCE ACTION",
		Short:   "Check if a subject has permission for an action on a resource",
		Example: checkExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(args))
			cmdutil.CheckErr(o.Run())
		},
	}

	return cmd
}

// Complete completes all the required options.
func (o *CheckOptions) Complete(args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("requires exactly 3 args: subject, resource, action")
	}
	o.Subject = args[0]
	o.Resource = args[1]
	o.Action = args[2]

	return nil
}

// Run executes a check subcommand using the specified options.
func (o *CheckOptions) Run() error {
	// Get authz server address from config, default to 9090
	authzServer := viper.GetString("authzserver.address")
	if authzServer == "" {
		authzServer = "http://127.0.0.1:9090"
	}

	reqBody := &v1.Request{
		Subject:  o.Subject,
		Resource: o.Resource,
		Action:   o.Action,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/authz", authzServer)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authz check failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the wrapped API response
	var apiResp api.Response
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !apiResp.Success {
		return fmt.Errorf("authz check failed: %s", apiResp.Msg)
	}

	// Extract the authz response from data
	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal authz data: %w", err)
	}

	var authzResp v1.Response
	if err := json.Unmarshal(dataBytes, &authzResp); err != nil {
		return fmt.Errorf("failed to parse authz response: %w", err)
	}

	if authzResp.Allowed {
		fmt.Fprintln(o.IOStreams.Out, "allowed")
	} else {
		fmt.Fprintf(o.IOStreams.Out, "denied: %s\n", authzResp.Reason)
	}

	return nil
}
