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

var (
	query  string
	format []string
	size   int
	index  string
	sort   []string
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Perform a search request.",
	Long:  `Perform a search request.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Start Search UI.
		ui.SearchDashboard(query, size, index, format, sort)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	// define flags.
	searchCmd.Flags().StringVarP(&query, "query", "q", "", "specify the query")
	searchCmd.Flags().StringSliceVarP(&sort, "sort", "s", []string{""}, "specify the sorting filter")
	searchCmd.Flags().IntVarP(&size, "size", "n", 15, "specify the size")
	searchCmd.Flags().StringSliceVarP(&format, "format", "f", []string{""}, "specify the format")
	searchCmd.Flags().StringVarP(&index, "index", "i", "*", "specify the index")

}
