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
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/spf13/viper"
)

//MainWindow creates the Main Menu Window.
type MainWindow struct {
	app        *tview.Application
	pages      *tview.Pages
	grid       *tview.Grid
	dashboards []string
	dashform   *tview.Form
	dashActive bool
	infoModal  *tview.Modal
	infoActive bool
}

// BANNER of the qelk.
const BANNER = `                    
                    ████  █████     
                   ░░███ ░░███      
  ████████  ██████  ░███  ░███ █████
 ███░░███  ███░░███ ░███  ░███░░███ 
░███ ░███ ░███████  ░███  ░██████░  
░███ ░███ ░███░░░   ░███  ░███░░███ 
░░███████ ░░██████  █████ ████ █████
 ░░░░░███  ░░░░░░  ░░░░░ ░░░░ ░░░░░ 
     ░███                           
     █████                          
	░░░░░                    v%s`

// VERSION of the qelk.
const VERSION = "0.1.0"

var (
	query        string
	format       []string
	size         int
	index        string
	sort         []string
	selectedDash string
)

// MainDashboard initiates the main menu.
func MainDashboard() {
	// Init. Main Dashboard
	w := &MainWindow{
		grid:  tview.NewGrid(),
		pages: tview.NewPages(),
	}
	// Init. Main App.
	w.app = w.mainapp()

	// Reads the config file and sets the dashboards.
	w.dashboards = setDashboards()

	// Init. Info modal.
	w.infoModal = w.infoM("Works!")

	// Init. text view of banner.
	header := createBanner()

	// Init. Stats button.
	mainStatButton := createMainButtons("Stats [F1]")
	mainStatButton.SetSelectedFunc(func() {
		w.app.Stop()
		StatsDashboard()
	})

	// Init. Dashboards button.
	mainDashButton := createMainButtons("Dashboards [F2]")
	mainDashButton.SetSelectedFunc(func() {
		w.pages.ShowPage("dashmodal")
	})

	// Init. Configuration Information text view.
	mainConfigGridCell := createConfigGrid()

	// Init. Author Information text view.
	mainAuthor := createAuthorGrid()

	// Init. Dashboard selection form.
	w.dashform = w.createDashForm()

	// Set grid configs and add widgets to grid.
	w.grid.SetRows(0, 5, 3)
	w.grid.SetColumns(0, 0)
	w.grid.AddItem(centeredPrimitives(header, 38, 13), 0, 0, 1, 2, 0, 0, false)
	w.grid.AddItem(mainStatButton, 1, 0, 1, 1, 0, 0, false)
	w.grid.AddItem(mainDashButton, 1, 1, 1, 1, 0, 0, false)
	w.grid.AddItem(mainConfigGridCell, 2, 0, 1, 1, 0, 0, false)
	w.grid.AddItem(mainAuthor, 2, 1, 1, 1, 0, 0, false)
	w.grid.SetBackgroundColor(tcell.ColorBlack)
	w.grid.SetBorder(true)
	w.grid.SetBorderColor(tcell.Color121)

	//Set pages for application.
	w.pages.AddPage("grid", w.grid, true, true)
	w.pages.AddPage("dashmodal", centeredPrimitives(w.dashform, 40, 10), true, false)
	w.pages.AddPage("info", centeredPrimitives(w.infoModal, 40, 7), true, false)

	// Start go routine.
	if err := w.app.SetRoot(w.pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

// Create app, set keys for app.
func (w *MainWindow) mainapp() *tview.Application {
	app := tview.NewApplication()
	defer app.Stop()
	app.SetBeforeDrawFunc(func(s tcell.Screen) bool {
		s.Clear()
		return false
	})
	// Set keys for app.
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			app.Stop()
			StatsDashboard()
			return nil
		}
		if event.Key() == tcell.KeyF2 {
			if w.dashActive == false {
				if viper.IsSet("dashboards") {
					w.pages.ShowPage("dashmodal")
					w.dashActive = true
				} else {
					w.infoModal.SetText("No dashboard found!")
					w.pages.ShowPage("info")
					w.infoActive = true
				}
			} else if w.dashActive == true {
				w.pages.HidePage("dashmodal")
			}

			return nil
		}
		return event
	})

	return app
}

// Get names of dashboards from config file.
func setDashboards() []string {
	adashboard := viper.GetStringMapString("dashboards")
	keys := []string{}
	for key := range adashboard {
		keys = append(keys, key)
	}

	return keys
}

// Create buttons for Main menu.
func createMainButtons(title string) *tview.Button {
	m := tview.NewButton(title)
	m.SetBorder(true)
	m.SetBorderColor(tcell.ColorYellow)
	m.SetBackgroundColor(tcell.ColorBlack)
	m.SetBackgroundColorActivated(tcell.ColorDarkCyan)

	return m
}

// Create config text view, set information about config file.
func createConfigGrid() *tview.TextView {
	m := tview.NewTextView()
	if err := viper.ReadInConfig(); err == nil {
		m.SetText("Using config file: " + viper.ConfigFileUsed())
	} else {
		m.SetText("Config file not found!")
	}
	m.SetTextAlign(tview.AlignLeft)
	m.SetTextColor(tcell.ColorYellow)
	m.SetBorder(true)
	m.SetBorderColor(tcell.ColorBlack)
	m.SetBackgroundColor(tcell.ColorBlack)

	return m
}

// Create author text view.
func createAuthorGrid() *tview.TextView {
	m := tview.NewTextView()
	m.SetText("github.com/serdarkkts")
	m.SetTextAlign(tview.AlignRight)
	m.SetTextColor(tcell.ColorYellow)
	m.SetBorder(true)
	m.SetBorderColor(tcell.ColorBlack)
	m.SetBackgroundColor(tcell.ColorBlack)

	return m
}

// Create banner text view.
func createBanner() *tview.TextView {
	h := tview.NewTextView()
	h.SetDynamicColors(true)
	h.SetText(fmt.Sprintf(BANNER, VERSION))
	h.SetTextColor(tcell.ColorYellow)
	h.SetTextAlign(tview.AlignCenter)
	return h
}

// Create dashboard selection form.
func (w *MainWindow) createDashForm() *tview.Form {
	f := tview.NewForm()
	// dashboards from config file.
	f.AddDropDown("Select a dashboard: ", w.dashboards, 0, func(optionLabel string, optionIndex int) {
		if optionLabel != "" {
			selectedDash = optionLabel
		}
	})
	// Run dashboard.
	f.AddButton("Inspect", func() {
		if viper.IsSet("dashboards." + selectedDash) {
			if viper.IsSet("dashboards." + selectedDash + ".index") {
				index = viper.GetString("dashboards." + selectedDash + ".index")
			}
			if viper.IsSet("dashboards." + selectedDash + ".size") {
				size = viper.GetInt("dashboards." + selectedDash + ".size")
			}
			if viper.IsSet("dashboards." + selectedDash + ".format") {
				format = viper.GetStringSlice("dashboards." + selectedDash + ".format")
			}
			if viper.IsSet("dashboards." + selectedDash + ".query") {
				query = viper.GetString("dashboards." + selectedDash + ".query")
			}
			if viper.IsSet("dashboards." + selectedDash + ".sort") {
				sort = viper.GetStringSlice("dashboards." + selectedDash + ".sort")
			}
			w.app.Stop()
			SearchDashboard(query, size, index, format, sort)
		} else {
			//
		}
	})
	f.AddButton("Close", func() {
		if w.dashActive == true {
			w.pages.HidePage("dashmodal")
			w.dashActive = false
		}
	})
	f.SetBorder(true).SetTitle("Dashboards").SetTitleAlign(tview.AlignCenter)
	f.SetButtonsAlign(tview.AlignCenter)

	return f
}

// Create information modal.
func (w *MainWindow) infoM(title string) *tview.Modal {
	m := tview.NewModal()
	m.SetText(title)
	m.AddButtons([]string{"Ok"})
	m.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Ok" {
			if w.infoActive == true {
				w.pages.HidePage("info")
				w.infoActive = false
			}
		}
	})
	return m
}
