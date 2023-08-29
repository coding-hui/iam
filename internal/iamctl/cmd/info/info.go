// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package info print the host information.
package info

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"

	cmdutil "github.com/coding-hui/iam/internal/iamctl/cmd/util"
	"github.com/coding-hui/iam/internal/iamctl/util/templates"
	"github.com/coding-hui/iam/pkg/cli/genericclioptions"

	"github.com/coding-hui/common/util/iputil"
)

// Info defines the host information struct.
type Info struct {
	HostName  string
	IPAddress string
	OSRelease string
	CPUCore   uint64
	MemTotal  string
	MemFree   string
}

// InfoOptions is an options struct to support 'info' sub command.
type InfoOptions struct {
	genericclioptions.IOStreams
}

var infoExample = templates.Examples(`
		# Print the host information
		iamctl info`)

// NewInfoOptions returns an initialized InfoOptions instance.
func NewInfoOptions(ioStreams genericclioptions.IOStreams) *InfoOptions {
	return &InfoOptions{
		IOStreams: ioStreams,
	}
}

// NewCmdInfo returns new initialized instance of 'info' sub command.
func NewCmdInfo(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewInfoOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "info",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Print the host information",
		Long:                  "Print the host information.",
		Example:               infoExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	return cmd
}

// Run executes an info sub command using the specified options.
func (o *InfoOptions) Run(args []string) error {
	var info Info

	memStat, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("get mem stat failed!error:%w", err)
	}

	info.MemTotal = strconv.FormatUint(memStat.Total, 10) + "M"
	info.MemFree = strconv.FormatUint(memStat.Free, 10) + "M"
	info.IPAddress = iputil.GetLocalIP()

	s := reflect.ValueOf(&info).Elem()
	typeOfInfo := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		v := fmt.Sprintf("%v", f.Interface())
		if v != "" {
			fmt.Fprintf(o.Out, "%12s %v\n", typeOfInfo.Field(i).Name+":", f.Interface())
		}
	}

	return nil
}
