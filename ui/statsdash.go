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
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/serdarkkts/qelk/es"
	"strings"
	"time"
)

// Refreshing stats
const refreshInterval = 500 * time.Millisecond

// StatsWindow  stats UI.
type StatsWindow struct {
	app             *tview.Application
	pages           *tview.Pages
	grid            *tview.Grid
	errmodal        *tview.Modal
	gridNodes       *tview.TextView
	gridHeader      *tview.TextView
	gridHealth      *tview.TextView
	gridFooter      *tview.TextView
	gridFilter      *tview.TextView
	gridIndices     *tview.TextView
	gridShards      *tview.TextView
	gridJvm         *tview.TextView
	gridDisk        *tview.TextView
	gridMem         *tview.TextView
	gridCPU         *tview.TextView
	gridVersions    *tview.TextView
	es              *elasticsearch.Client
	nodeFilter      tview.Primitive
	nodeFilterConst string
	filterPre       string
	filterIP        string
	filterActive    bool
}

var (
	conerr      error
	json        string
	nodefilters []string = []string{"_all", "master:true", "data:true", "ingest:true", "voting_only:true", "ml:true",
		"coordinating_only:true", "master:false", "data:false", "ingest:false", "voting_only:false", "ml:false", "coordinating_only:false"}
)

//StatsDashboard initiates the Stats window.
func StatsDashboard() {

	// Init. Stats Dashboard
	w := &StatsWindow{
		gridHeader:      statsHeader(""),
		gridHealth:      statsHealth(""),
		gridFooter:      statsFooter(),
		gridFilter:      statsFilter(""),
		gridNodes:       statsgridcell("Total Nodes: "+"\n[green]Successful: "+"\n[red]Failed: ", " Nodes "),
		gridIndices:     statsgridcell("Total: "+"\nDocuments: "+"\nUptime: ", " Indices "),
		gridShards:      statsgridcell("Total: "+"\nPrimary: "+"\nReplication: "+"\nSize: ", " Shards "),
		gridJvm:         statsgridcell("Used: "+"\nMax: "+"\nThreads: ", " JVM "),
		gridDisk:        statsgridcell("Available: "+"\nFree: "+"\nTotal: ", " Disk "),
		gridMem:         statsgridcell("Free: "+"\nUsed: "+"\nTotal: ", " Memory | "+"[red]%"+" "),
		gridCPU:         statsgridcell("Available Proc: "+"\nAllocated Proc: ", " CPU | "+"[red]%"+" "),
		gridVersions:    statsgridcell("", " Version "),
		es:              esc.Connect(),
		nodeFilterConst: "_all",
		filterPre:       "_all",
	}

	// Init. App
	w.app = w.statsApp()

	// Init. Filter Modal
	w.nodeFilter = centeredPrimitives(w.nodeFilterForm(), 40, 16)

	// Init. node filter form
	w.errmodal = w.errModal()

	// Init. main grid
	w.grid = w.mainGrid()

	// Init. pages for ui.
	w.pages = statsPages(w.grid, w.errmodal, w.nodeFilter)

	// Start setting values to gridcells concurrently.
	go w.updateStats()

	// Start go routine.
	if err := w.app.SetRoot(w.pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

// Setting values to gridcells concurrently.
func (w *StatsWindow) updateStats() {
	for {
		time.Sleep(refreshInterval)
		json, conerr = esc.Stats(w.es, w.nodeFilterConst)
		if conerr == nil {
			w.app.QueueUpdateDraw(func() {
				w.gridHeader.SetText(strings.ToUpper(v("cluster_name")))
				w.gridHealth.SetText(strings.ToUpper(v("status")))
				w.gridHealth.SetTextColor(statsHealthColor(strings.ToUpper(v("status"))))
				w.gridHealth.SetBorderColor(statsHealthColor(strings.ToUpper(v("status"))))
				w.gridNodes.SetText("Total Nodes: " + v("_nodes.total") + "\n[green]Successful: " + v("_nodes.successful") + "\n[red]Failed: " + v("_nodes.failed"))
				w.gridIndices.SetText("Total: " + v("indices.count") + "\nDocuments: " + v("indices.docs.count") + "\nUptime: " + v("nodes.jvm.max_uptime"))
				w.gridShards.SetText("Total: " + v("indices.shards.total") + "\nPrimary: " + v("indices.shards.primaries") + "\nReplication: " + v("indices.shards.replication") + "\nSize: " + v("indices.store.size"))
				w.gridJvm.SetText("Used: " + v("nodes.jvm.mem.heap_used") + "\nMax: " + v("nodes.jvm.mem.heap_max") + "\nThreads: " + v("nodes.jvm.threads"))
				w.gridDisk.SetText("Available: " + v("nodes.fs.available") + "\nFree: " + v("nodes.fs.free") + "\nTotal: " + v("nodes.fs.total"))
				w.gridMem.SetText("Free: " + v("nodes.os.mem.free") + "\nUsed: " + v("nodes.os.mem.used") + "\nTotal: " + v("nodes.os.mem.total"))
				w.gridMem.SetTitle(" Memory | " + "[red]%" + v("nodes.os.mem.used_percent") + " ")
				w.gridCPU.SetText("Available Proc: " + v("nodes.os.available_processors") + "\nAllocated Proc: " + v("nodes.os.allocated_processors"))
				w.gridCPU.SetTitle(" CPU | " + "[red]%" + v("nodes.process.cpu.percent") + " ")
				w.gridVersions.SetText(v("nodes.versions"))
				w.gridFilter.SetText(w.nodeFilterConst)
			})
		} else if conerr != nil {
			w.app.QueueUpdateDraw(func() {
				w.pages.ShowPage("modal")
			})
		}

	}
}

// mainGrid creates the main grid for the Stats Window.
func (w *StatsWindow) mainGrid() *tview.Grid {
	g := tview.NewGrid()
	g.SetRows(3, 0, 0, 3)
	g.SetColumns(0, 0, 0, 0)
	g.AddItem(w.gridHeader, 0, 0, 1, 3, 0, 0, false)
	g.AddItem(w.gridHealth, 0, 3, 1, 1, 0, 0, false)
	g.AddItem(w.gridCPU, 1, 0, 1, 1, 0, 0, false)
	g.AddItem(w.gridIndices, 2, 0, 1, 1, 0, 0, false)
	g.AddItem(w.gridMem, 1, 1, 1, 1, 0, 0, false)
	g.AddItem(w.gridJvm, 2, 1, 1, 1, 0, 0, false)
	g.AddItem(w.gridDisk, 1, 2, 1, 1, 0, 0, false)
	g.AddItem(w.gridShards, 2, 2, 1, 1, 0, 0, false)
	g.AddItem(w.gridNodes, 1, 3, 1, 1, 0, 0, false)
	g.AddItem(w.gridVersions, 2, 3, 1, 1, 0, 0, false)
	g.AddItem(w.gridFooter, 3, 0, 1, 3, 0, 0, false)
	g.AddItem(w.gridFilter, 3, 3, 1, 1, 0, 0, false)

	return g
}

// errModal Init. Error modal for Stats Window.
func (w *StatsWindow) errModal() *tview.Modal {
	m := tview.NewModal()
	m.SetText("Connection failed. Trying again in few seconds.")
	m.AddButtons([]string{"Quit", "Cancel"})
	m.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Quit" {
			w.app.Stop()
		}
		if buttonLabel == "Cancel" {
			w.pages.HidePage("modal")
		}
	})
	return m
}

// nodeFilterForm takes filter inputs from user and applies them to Stats.
func (w *StatsWindow) nodeFilterForm() *tview.Form {
	f := tview.NewForm()
	f.AddDropDown("Select a filter: ", nodefilters, 0, func(optionLabel string, optionIndex int) {
		if optionLabel != "" {
			w.filterPre = optionLabel
		}
	})
	f.AddInputField("Node IP: ", "", 0, nil, func(inputIp string) {
		if inputIp != "" {
			w.filterIP = inputIp
		}
	})
	f.AddButton("Apply", func() {
		if len(w.filterIP) > 1 {
			w.nodeFilterConst = w.filterIP
		} else {
			w.nodeFilterConst = w.filterPre
		}
	})
	f.AddButton("Close", func() {
		w.pages.HidePage("filter")
		w.filterActive = false
	})
	f.SetBorder(true).SetTitle(" Define a Node Filter! ").SetTitleAlign(tview.AlignCenter)
	f.SetButtonsAlign(tview.AlignCenter)
	return f
}

// statsApp Init. App for Stats Window.
func (w *StatsWindow) statsApp() *tview.Application {

	a := tview.NewApplication()

	// tcell ColorDefault issue: https://github.com/rivo/tview/issues/270#issuecomment-485083503
	defer a.Stop()
	a.SetBeforeDrawFunc(func(s tcell.Screen) bool {
		s.Clear()
		return false
	})

	// Init. Keys for application
	a.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			if w.filterActive == false {
				w.pages.ShowPage("filter")
				w.filterActive = true
			} else if w.filterActive == true {
				w.pages.HidePage("filter")
				w.filterActive = false
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

	return a

}
