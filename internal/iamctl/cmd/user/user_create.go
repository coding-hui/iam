// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	metav1 "github.com/coding-hui/common/meta/v1"
	apiclientv1 "github.com/coding-hui/wecoding-sdk-go/wecoding/iam/apiserver/v1"

	cmdutil "github.com/coding-hui/iam/internal/iamctl/cmd/util"
	"github.com/coding-hui/iam/internal/iamctl/util/templates"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/cli/genericclioptions"
)

const (
	createUsageStr = "create USERNAME PASSWORD EMAIL"
)

// CreateOptions is an options struct to support create subcommands.
type CreateOptions struct {
	Alias string
	Phone string

	User *v1.CreateUserRequest

	Client apiclientv1.APIV1Interface
	genericclioptions.IOStreams
}

var (
	createLong = templates.LongDesc(`Create a user on iam platform.
If nickname not specified, username will be used.`)

	createExample = templates.Examples(`
		# Create user with given input
		iamctl user create wecoding wecoding@2023 wecoding@yeah.net

		# Create user with phone and alias
		iamctl user create wecoding wecoding@2023 wecoding@yeah.net --phone=1897777xxx --alias=coder`)

	createUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nUSERNAME, PASSWORD and EMAIL are required arguments for the create command",
		createUsageStr,
	)
)

// NewCreateOptions returns an initialized CreateOptions instance.
func NewCreateOptions(ioStreams genericclioptions.IOStreams) *CreateOptions {
	return &CreateOptions{
		IOStreams: ioStreams,
	}
}

// NewCmdCreate returns new initialized instance of create sub command.
func NewCmdCreate(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewCreateOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   createUsageStr,
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Create a user resource",
		TraverseChildren:      true,
		Long:                  createLong,
		Example:               createExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	// mark flag as deprecated
	cmd.Flags().StringVar(&o.Alias, "alias", o.Alias, "The alias of the user.")
	cmd.Flags().StringVar(&o.Phone, "phone", o.Phone, "The phone number of the user.")

	return cmd
}

// Complete completes all the required options.
func (o *CreateOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	if len(args) < 3 {
		return cmdutil.UsageErrorf(cmd, createUsageErrStr)
	}

	if o.Alias == "" {
		o.Alias = args[0]
	}

	o.User = &v1.CreateUserRequest{
		Name:     args[0],
		Password: args[1],
		Email:    args[2],
		Alias:    o.Alias,
		Phone:    o.Phone,
	}

	clientConfig, err := f.ToRESTConfig()
	if err != nil {
		return err
	}
	o.Client, err = apiclientv1.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *CreateOptions) Validate(cmd *cobra.Command, args []string) error {
	if errs := o.User.Validate(); len(errs) != 0 {
		return errs.ToAggregate()
	}

	return nil
}

// Run executes a create subcommand using the specified options.
func (o *CreateOptions) Run(args []string) error {
	ret, err := o.Client.Users().Create(context.TODO(), o.User, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "user/%s created\n", ret.Name)

	return nil
}
