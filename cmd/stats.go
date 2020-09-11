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
	"github.com/serdarkkts/qelk/ui"
	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Monitor your Elasticsearch cluster.",
	Long:  `Helps you to monitor your Elasticsearch cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		// Start Stats UI.
		ui.StatsDashboard()
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
