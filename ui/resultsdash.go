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
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strconv"
)

//SearchWindow search UI.
type SearchWindow struct {
	app           *tview.Application
	pages         *tview.Pages
	grid          *tview.Grid
	footer        *tview.TextView
	previewActive bool
	preview       *tview.Modal
	sResults      *Results
	totalResults  *tview.TextView
}

//SearchDashboard  initiates the Search window.
func SearchDashboard(query string, size int, index string, fields []string, sort []string) {

	// Init. Search Dashboard
	w := &SearchWindow{
		grid:   tview.NewGrid(),
		footer: setFooter(),
	}
	// Init. App
	w.app = w.apps()

	// Get results.
	w.sResults = SearchResults(size, index, query, fields, sort)

	// Init. Inspect modal for long results.
	w.preview = w.sResults.InspectResult()

	// Init. total results.
	w.totalResults = w.setTotalResults()

	// Init. flew rows for results.
	f := tview.NewFlex()
	f.SetDirection(tview.FlexRow)
	f.SetBackgroundColor(tcell.ColorBlack)
	f.SetBorder(false)
	f.SetBorderColor(tcell.Color220)
	f.SetTitle("Logs")
	f.AddItem(w.sResults, 0, 1, true)

	// Config Search window grid.
	w.grid.SetRows(0, 3)
	w.grid.SetColumns(-5, -1)
	w.grid.AddItem(f, 0, 0, 1, 2, 0, 0, true)
	w.grid.AddItem(w.footer, 1, 0, 1, 1, 0, 0, false)
	w.grid.AddItem(w.totalResults, 1, 1, 1, 1, 0, 0, false)
	w.grid.SetBackgroundColor(tcell.ColorBlack)

	// Init. pages
	w.pages = tview.NewPages().
		AddPage("background", w.grid, true, true).
		AddPage("modal", centeredPrimitives(w.preview, 5, 5), false, false)

	// Start go routine.
	if err := w.app.SetRoot(w.pages, true).Run(); err != nil {
		panic(err)
	}

}

// Set navbar of Search window.
func setFooter() *tview.TextView {
	h := tview.NewTextView()
	h.SetDynamicColors(true)
	h.SetText("[darkcyan]ENTER:[white]Inspect |[darkcyan] F2: [white]Quit |[darkcyan] F3: [white]Main Menu |[darkcyan] Use arrow keys to scroll. ")
	h.SetTitleAlign(tview.AlignLeft)
	h.SetBackgroundColor(tcell.ColorBlack)
	h.SetBorder(true)
	h.SetBorderColor(tcell.Color66)
	return h
}

// Set total results.
func (w *SearchWindow) setTotalResults() *tview.TextView {
	h := tview.NewTextView()
	h.SetDynamicColors(true)
	h.SetText(strconv.Itoa(len(w.sResults.cells)))
	h.SetBackgroundColor(tcell.ColorBlack)
	h.SetTitle("Total")
	h.SetTitleAlign(tview.AlignLeft)
	h.SetBorder(true)
	h.SetBorderColor(tcell.Color220)
	return h
}

// Init app for SearchWindow.
func (w *SearchWindow) apps() *tview.Application {

	app := tview.NewApplication()
	defer app.Stop()
	app.SetBeforeDrawFunc(func(s tcell.Screen) bool {
		s.Clear()
		return false
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			if w.previewActive == false {
				w.preview.SetText(w.sResults.cells[w.sResults.focus].GetText(true))
				w.pages.ShowPage("modal")
				w.previewActive = true
			} else if w.previewActive == true {
				w.pages.HidePage("modal")
				w.previewActive = false
			}
			return nil
		}
		if event.Key() == tcell.KeyF2 {
			w.app.Stop()
			return nil
		}
		if event.Key() == tcell.KeyF3 {
			w.app.Stop()
			MainDashboard()
			return nil
		}
		return event
	})

	return app
}
