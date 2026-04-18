// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/coding-hui/iam/internal/apiserver"
	"github.com/coding-hui/iam/internal/config"
)

func main() {
	// Check for explicit config file
	configFile := os.Getenv("IAM_CONFIG_FILE")
	for i, arg := range os.Args[1:] {
		if arg == "--config" && i+1 < len(os.Args)-1 {
			configFile = os.Args[i+2]
			break
		}
		if strings.HasPrefix(arg, "--config=") {
			configFile = strings.TrimPrefix(arg, "--config=")
			break
		}
	}

	viper.SetConfigName("apiserver")
	viper.SetConfigType("yaml")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".iam", "conf"))
		viper.AddConfigPath("/etc/iam")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read config: %v\n", err)
		os.Exit(1)
	}

	zap.S().Infof("Config file used: %s", viper.ConfigFileUsed())

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to unmarshal config: %v\n", err)
		os.Exit(1)
	}

	// Validate
	if cfg.Database.DSN == "" {
		fmt.Fprintf(os.Stderr, "Error: database DSN is required\n")
		os.Exit(1)
	}
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		fmt.Fprintf(os.Stderr, "Error: server port must be between 1 and 65535\n")
		os.Exit(1)
	}

	if err := apiserver.Run("apiserver", &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
