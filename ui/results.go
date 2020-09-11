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

package ui

import (
	"fmt"
	"github.com/serdarkkts/qelk/es"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/tidwall/gjson"
	"strings"
)

// Line size per document.
const templateLines = 5

// Results implements tview's Grid for search response.
type Results struct {
	*tview.Grid
	cells         []*tview.TextView
	focus         int
	customHandler func(*tcell.EventKey) *tcell.EventKey
	es            *elasticsearch.Client
	fields        []string
}

// SearchResults initiates search results window.
func SearchResults(sizes int, index string, query string, fields []string, sort []string) *Results {
	// Init. grid.
	var grid = tview.NewGrid()
	// For fields of results.
	var sb strings.Builder

	// Init. Results.
	c := &Results{
		Grid: grid,
		es:   esc.Connect(),
	}

	// Set grid config.
	grid.SetColumns(0)
	grid.SetBorders(false)
	grid.SetInputCapture(c.DefaultInputHandler)
	grid.SetGap(1, 1)
	// Perform search request.
	jsonresults, _ := esc.Search(c.es, sizes, index, query, fields, sort)

	// Reading only source field. This ignores the indices of results.(must check)
	r := gjson.Get(jsonresults, "hits.hits.#._source")

	// Create rows and array for results.
	var rows = make([]int, len(r.Array()))
	c.cells = make([]*tview.TextView, len(r.Array()))

	// Check if at least one field provided.
	if fields[0] == "" {
		c.fields = []string{"@this"}
	} else {
		c.fields = fields
	}

	// Set results to the text views.
	for i, hit := range r.Array() {
		// Set line size per result.
		rows[i] = templateLines

		// A simple trick for making fields more visible, change color one of every two fields.
		for j := range c.fields {
			fieldcolor := "[darkviolet]"
			if j%2 == 0 {
				fieldcolor = "[white]"
			}
			// Create format of results.
			fmt.Fprintf(&sb, "%s%s ", fieldcolor, hit.Get(c.fields[j]).String())
		}
		// Config of result text views.
		resultCell := tview.NewTextView()
		resultCell.SetScrollable(true)
		resultCell.SetDynamicColors(true)

		// Set result to text view.
		resultCell.SetText(sb.String())

		// Add every text view to grid.
		grid.AddItem(resultCell, i, 0, 1, 1, 0, 0, false)

		// Set results to cells array for scrolling.
		c.cells[i] = resultCell

		// Reset string builder
		sb.Reset()
	}
	// Set line sizes of grid.
	grid.SetRows(rows...)

	return c
}

// DefaultInputHandler is for setting keys of results page.
func (c *Results) DefaultInputHandler(ev *tcell.EventKey) *tcell.EventKey {
	old := c.focus
	switch ev.Key() {
	// dec. focus variable for scroll
	case tcell.KeyUp:
		if c.focus == 0 {
			return ev
		}

		c.focus--
	// inc. focus
	case tcell.KeyDown:
		if c.focus == len(c.cells)-1 {
			return ev
		}

		c.focus++

	default:
		if c.customHandler != nil {
			return c.customHandler(ev)
		}

		return ev
	}

	// for making scroll operation more visible, change colors of rows.
	if len(c.cells) > 1 {

		c.cells[old].SetBackgroundColor(tcell.ColorBlack)
		c.cells[old].SetTextColor(tcell.ColorWhite)

		c.cells[c.focus].SetBackgroundColor(tcell.ColorGrey)
		c.cells[c.focus].SetTextColor(tcell.ColorBlack)

		c.SetOffset(c.focus-1, 0)
	}

	return nil
}

// InspectResult creates a modal for inspecting long results.
func (c *Results) InspectResult() *tview.Modal {
	m := tview.NewModal()
	m.SetText("Nothing yet")
	m.AddButtons([]string{"Ok"})
	m.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Ok" {
			// c.app.Stop()
		}
	})
	return m
}
