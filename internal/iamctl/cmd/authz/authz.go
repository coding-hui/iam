// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package authz provides commands for authorization checks via iam-authz-server.
package authz

import (
	"github.com/spf13/cobra"

	cmdutil "github.com/coding-hui/iam/internal/iamctl/cmd/util"
	"github.com/coding-hui/iam/internal/iamctl/util/templates"
	"github.com/coding-hui/iam/pkg/cli/genericclioptions"
)

var authzLong = templates.LongDesc(`
	Authorization verification commands.

	Checks permissions against the iam-authz-server.`)

// NewCmdAuthz returns new initialized instance of 'authz' sub command.
func NewCmdAuthz(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "authz SUBCOMMAND",
		Short: "Authorization verification commands",
		Long:  authzLong,
		Run:   cmdutil.DefaultSubCommandRun(ioStreams.ErrOut),
	}

	cmd.AddCommand(NewCmdCheck(f, ioStreams))

	return cmd
}
