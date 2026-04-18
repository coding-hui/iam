// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/coding-hui/iam/internal/apiserver"
)

func main() {
	// Set default config path
	viper.SetConfigName("apiserver")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".iam", "conf"))
	viper.AddConfigPath("/etc/iam")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read config: %v\n", err)
		os.Exit(1)
	}

	zap.S().Infof("Config file used: %s", viper.ConfigFileUsed())

	opts := apiserver.NewOptions()
	if err := viper.Unmarshal(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to unmarshal config: %v\n", err)
		os.Exit(1)
	}

	if err := opts.Complete(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to complete config: %v\n", err)
		os.Exit(1)
	}

	if errs := opts.Validate(); len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}

	if err := apiserver.Run("apiserver"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
