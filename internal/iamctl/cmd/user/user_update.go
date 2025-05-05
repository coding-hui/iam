// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	metav1 "github.com/coding-hui/common/meta/v1"
	"github.com/coding-hui/wecoding-sdk-go/services/iam"

	cmdutil "github.com/coding-hui/iam/internal/iamctl/cmd/util"
	"github.com/coding-hui/iam/internal/iamctl/util/templates"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/cli/genericclioptions"
)

const (
	updateUsageStr = "update USERNAME"
)

// UpdateOptions is an options struct to support update subcommands.
type UpdateOptions struct {
	ID    string
	Name  string
	Alias string
	Email string
	Phone string

	iamclient iam.IamInterface
	genericclioptions.IOStreams
}

var (
	updateLong = templates.LongDesc(`Update a user resource. 

Can only update alias, email and phone.

NOTICE: field will be updated to zero value if not specified.`)

	updateExample = templates.Examples(`
		# Update use foo's information
		iamctl user update foo --alias=foo2 --email=foo@qq.com --phone=1812883xxxx`)

	updateUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nUSERID is required arguments for the update command",
		updateUsageStr,
	)
)

// NewUpdateOptions returns an initialized UpdateOptions instance.
func NewUpdateOptions(ioStreams genericclioptions.IOStreams) *UpdateOptions {
	return &UpdateOptions{
		IOStreams: ioStreams,
	}
}

// NewCmdUpdate returns new initialized instance of update sub command.
func NewCmdUpdate(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewUpdateOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   updateUsageStr,
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Update a user resource",
		TraverseChildren:      true,
		Long:                  updateLong,
		Example:               updateExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.Alias, "alias", o.Alias, "The alias of the user.")
	cmd.Flags().StringVar(&o.Email, "email", o.Email, "The email of the user.")
	cmd.Flags().StringVar(&o.Phone, "phone", o.Phone, "The phone number of the user.")

	return cmd
}

// Complete completes all the required options.
func (o *UpdateOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	if len(args) == 0 {
		return cmdutil.UsageErrorf(cmd, "%s", updateUsageErrStr)
	}

	o.ID = args[0]
	o.iamclient, err = f.IAMClient()
	if err != nil {
		return err
	}

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *UpdateOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes an update subcommand using the specified options.
func (o *UpdateOptions) Run(args []string) error {
	user, err := o.iamclient.APIV1().Users().Get(context.TODO(), o.ID, metav1.GetOptions{})
	if err != nil {
		return err
	}

	updateReq := &v1.UpdateUserRequest{
		Alias: o.Alias,
		Email: o.Email,
		Phone: o.Phone,
	}

	ret, err := o.iamclient.APIV1().Users().Update(context.TODO(), user.InstanceID, updateReq, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "user/%s updated\n", ret.Name)

	return nil
}
