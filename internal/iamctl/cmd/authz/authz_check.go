// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authz

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	v1 "github.com/coding-hui/iam/pkg/api/authzserver/v1"
	authzv1 "github.com/coding-hui/wecoding-sdk-go/services/iam/authz/v1"
	cmdutil "github.com/coding-hui/iam/internal/iamctl/cmd/util"
	"github.com/coding-hui/iam/internal/iamctl/util/templates"
	"github.com/coding-hui/iam/pkg/cli/genericclioptions"
)

// CheckOptions is an options struct to support authz check subcommand.
type CheckOptions struct {
	Subject  string
	Resource string
	Action   string

	authzClient authzv1.AuthzV1Interface
	genericclioptions.IOStreams
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
			cmdutil.CheckErr(o.Complete(f, args))
			cmdutil.CheckErr(o.Run())
		},
	}

	return cmd
}

// Complete completes all the required options.
func (o *CheckOptions) Complete(f cmdutil.Factory, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("requires exactly 3 args: subject, resource, action")
	}
	o.Subject = args[0]
	o.Resource = args[1]
	o.Action = args[2]

	var err error
	o.authzClient, err = f.AuthzV1Client()
	if err != nil {
		return err
	}

	return nil
}

// Run executes a check subcommand using the specified options.
func (o *CheckOptions) Run() error {
	req := &v1.Request{
		Subject:  o.Subject,
		Resource: o.Resource,
		Action:   o.Action,
	}

	resp, err := o.authzClient.Authz().Authorize(context.TODO(), req)
	if err != nil {
		return err
	}

	if resp.Allowed {
		fmt.Fprintln(o.Out, "allowed")
	} else {
		fmt.Fprintf(o.Out, "denied: %s\n", resp.Reason)
	}

	return nil
}
