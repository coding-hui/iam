// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package login

import (
	"github.com/spf13/cobra"

	"github.com/coding-hui/iam/internal/iamctl/cmd/util"
	"github.com/coding-hui/iam/pkg/cli/genericclioptions"
)

// NewCmdLogin creates a new login command.
func NewCmdLogin(f util.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to iam platform",
		Long:  "Login to iam platform using different authentication methods",
	}

	cmd.AddCommand(
		NewCmdDeviceLogin(f, ioStreams),
		// Can add other login methods here
		// NewCmdTokenLogin(f, ioStreams),
		// NewCmdPasswordLogin(f, ioStreams),
	)

	return cmd
}
