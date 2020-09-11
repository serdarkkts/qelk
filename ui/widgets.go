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
	"bytes"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/tidwall/gjson"
	"io"
)

var (
	statuscolor tcell.Color
	themeColor  tcell.Color = tcell.ColorBlack
)

// centeredPrimitives for centering widgets.
func centeredPrimitives(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

// statsPages Init. pages
func statsPages(grid *tview.Grid, modal *tview.Modal, filter tview.Primitive) *tview.Pages {
	p := tview.NewPages().
		AddPage("background", grid, true, true).
		AddPage("modal", modal, true, false).
		AddPage("filter", filter, true, false)
	return p
}

// Init. and sets the text and title fields of stat text views.
func statsgridcell(text string, title string) *tview.TextView {
	s := tview.NewTextView()
	s.SetDynamicColors(true)
	s.SetText(text)
	s.SetTextAlign(tview.AlignCenter)
	s.SetTitle(title)
	s.SetTitleAlign(tview.AlignCenter)
	s.SetTitleColor(tcell.Color110)
	s.SetBorder(true)
	s.SetBackgroundColor(themeColor)
	s.SetBorderPadding(1, 0, 0, 0)
	return s
}

// statsHeader Init. and sets the text field of Header.
func statsHeader(clusterName string) *tview.TextView {
	h := tview.NewTextView()
	h.SetText(clusterName)
	h.SetTextAlign(tview.AlignCenter)
	h.SetBorder(true)
	h.SetBorderAttributes(tcell.AttrBold)
	h.SetTitle(" CLUSTER ")
	h.SetTitleAlign(tview.AlignLeft)
	h.SetBackgroundColor(themeColor)
	return h
}

// statsHealth Init. and sets the color and text field of Status.
func statsHealth(clusterStatus string) *tview.TextView {
	h2 := tview.NewTextView()
	h2.SetText(clusterStatus)
	h2.SetTextAlign(tview.AlignCenter)
	h2.SetTextColor(tcell.ColorYellow)
	h2.SetBorder(true)
	h2.SetBorderAttributes(tcell.AttrUnderline)
	h2.SetTitle(" STATUS ")
	h2.SetTitleAlign(tview.AlignLeft)
	h2.SetBorderColor(tcell.ColorYellow)
	h2.SetBackgroundColor(themeColor)
	return h2
}

// statsFooter Init. and sets the text field of Footer.
func statsFooter() *tview.TextView {
	f := tview.NewTextView()
	f.SetText("[darkcyan]F1[white] Node Filter |[darkcyan] F2 [white]Quit |[darkcyan] F3 [white]Main Menu ")
	f.SetDynamicColors(true)
	f.SetBackgroundColor(themeColor)
	f.SetBorder(true)
	return f
}

// statsFilter Init. and sets the text field of Filter.
func statsFilter(filter string) *tview.TextView {
	f := tview.NewTextView()
	f.SetText(filter).SetTextAlign(tview.AlignCenter)
	f.SetTitle(" Filter ").SetTitleAlign(tview.AlignLeft)
	f.SetDynamicColors(true)
	f.SetBackgroundColor(themeColor)
	f.SetBorder(true)
	return f
}

// statsHealthColor dynamically sets the color of Status cell.
func statsHealthColor(clusterStatus string) tcell.Color {
	switch clusterStatus {
	case "GREEN":
		statuscolor = tcell.Color40
	case "YELLOW":
		statuscolor = tcell.ColorYellow
	case "RED":
		statuscolor = tcell.ColorRed
	default:
		statuscolor = tcell.ColorRed
	}
	return statuscolor
}

// v returns the json value.
func v(path string) string {
	value := gjson.Get(json, path).String()
	return value
}

// Read the response.
func Read(r io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.String()
}
