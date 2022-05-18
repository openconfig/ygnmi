// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/openconfig/ygnmi/app/ygnmi/cmd/generator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// New creates a new root command.
//nolint:errcheck
func New() *cobra.Command {
	root := &cobra.Command{
		Use: "ygnmi",
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		ygotVersion := "unknown"
		for _, dep := range info.Deps {
			if dep.Path == "github.com/openconfig/ygot" {
				ygotVersion = dep.Version
			}
		}
		root.Version = fmt.Sprintf("%s: (ygot: %s)", info.Main.Version, ygotVersion)
	}

	cfgFile := root.PersistentFlags().String("config_file", "", "Path to config file.")
	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if *cfgFile != "" {
			viper.SetConfigFile(*cfgFile)
			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("error reading config: %w", err)
			}
		}
		viper.BindPFlags(cmd.Flags())
		viper.AutomaticEnv()
		// Workaround for a flag marked as required, but set by config file or env var.
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if viper.IsSet(f.Name) {
				cmd.Flags().Set(f.Name, fmt.Sprint(viper.Get(f.Name)))
			}
		})
		return nil
	}

	root.AddCommand(generator.New())
	return root
}
