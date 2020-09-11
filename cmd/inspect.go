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
	"github.com/serdarkkts/qelk/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect your custom dashboards.",
	Long:  `Inspect your custom dashboards.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dashtemplate := strings.Join(args, " ")
		if viper.IsSet("dashboards." + dashtemplate) {
			if viper.IsSet("dashboards." + dashtemplate + ".index") {
				index = viper.GetString("dashboards." + dashtemplate + ".index")
			}
			if viper.IsSet("dashboards." + dashtemplate + ".size") {
				size = viper.GetInt("dashboards." + dashtemplate + ".size")
			}
			if viper.IsSet("dashboards." + dashtemplate + ".format") {
				format = viper.GetStringSlice("dashboards." + dashtemplate + ".format")
			}
			if viper.IsSet("dashboards." + dashtemplate + ".query") {
				query = viper.GetString("dashboards." + dashtemplate + ".query")
			}
			if viper.IsSet("dashboards." + dashtemplate + ".sort") {
				sort = viper.GetStringSlice("dashboards." + dashtemplate + ".sort")
			}
			ui.SearchDashboard(query, size, index, format, sort)
		} else {
			fmt.Println("Dashboard not found!")
		}

	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
