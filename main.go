// Copyright 2015 The Tops'l Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	//"time"
	//"io/ioutil"
	//"strconv"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

//var seshKeys = [8]rune{'a','s','d','f','j','k','l',';'}
//var seshColors = [8]tcell.Color{8,9,10,11,12,13,14,15}
//var seshColors = [8]tcell.Color{0,1,2,3,4,5,6,7}

type boxL struct {
	views.BoxLayout
}

var app = &views.Application{}
var box = &boxL{}
var seshBox = NewSeshBox()
var seshStatus = &boxL{}
var seshShell = &boxL{}
var buttons [8]*SeshButton
var seshLine = views.NewBoxLayout(views.Horizontal)

func (m *boxL) HandleEvent(ev tcell.Event) bool {
	//fmt.Fprintf(os.Stderr, "event\n")
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape, tcell.KeyEnter:
			app.Quit()
			return true
		}
	}
	return m.BoxLayout.HandleEvent(ev)
}

func (m *boxL) Refresh() {

}

func main() {
	//seshLine = views.NewBoxLayout(views.Horizontal)
/*
	//read in current dir
	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var spacerText = views.NewText()
	spacerText.SetText(" ")
	seshLine.AddWidget(spacerText,0)
	// create sesh buttons, add dirname to them
	for i := 0; i < 8; i++ {
		buttons[i] = NewButton()
		buttons[i].SetStyle(tcell.StyleDefault.Foreground(seshColors[i]).Background(tcell.ColorGrey))
		buttons[i].SetAlignment(views.AlignMiddle)
		buttons[i].Key = seshKeys[i]

		if(len(files) > i) {
			buttons[i].SetFileInfo(files[i])
		}
		seshLine.AddWidget(buttons[i], 0.1)

		spacerText = views.NewText()
		spacerText.SetText(" ")
		if i == 3 {
			spacerText.SetText("  ") // bigger gap between hands
		}
		seshLine.AddWidget(spacerText, 0)
	}

	seshBox.Initialize()

	seshBox.AddWidget(seshStatus, 0.5)
	seshBox.AddWidget(seshLine, 0)

	seshBox.AddWidget(seshShell, 0.5)
*/

	seshBox.Initialize()
	seshBox.SetOrientation(views.Vertical)

	app.SetRootWidget(seshBox)

	seshBox.App = app
	if e := app.Run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}
}
