// Copyright © 2020 Serdar KÖKTAŞ <contact@serdarkoktas.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/serdarkkts/qelk/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qelk",
	Short: "A terminal UI to monitor and query Elasticsearch.",
	Long: `qelk helps you to monitor your Elasticsearch cluster, execute search queries
and creates dashboards for your search queries.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Start Main Menu UI.
		ui.MainDashboard()
	},
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Flags and config
func init() {
	log.SetFlags(0)
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.qelk.yaml)")
}

// initConfig reads in config file if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".qelk" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".qelk")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//..
	}
}
