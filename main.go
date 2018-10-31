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
	"io/ioutil"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

var seshKeys = [8]rune{'a','s','d','f','j','k','l',';'}
var seshColors = [8]tcell.Color{15,9,10,11,12,13,14,15}

type boxL struct {
	views.BoxLayout
}

var app = &views.Application{}
var box = &boxL{}
//var SeshButtons [8]*views.Text
var buttons [8]*SeshButton


func (m *boxL) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape {
			app.Quit()
			return true
		}
		switch ev.Key() {
		case tcell.KeyEscape, tcell.KeyEnter: 
			app.Quit()
			return true
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'a':
			case 's':
			case 'd':
			}
		}
	}
	return m.BoxLayout.HandleEvent(ev)
}

func main() {
	seshLine := views.NewBoxLayout(views.Horizontal)

	//read in current dir
	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// create sesh buttons, add dirname to them
	for i := 0; i < 8; i++ {
		buttons[i] = NewButton()
		buttons[i].SetStyle(tcell.StyleDefault.Foreground(seshColors[i]))
		buttons[i].Key = seshKeys[i]

		buttons[i].CalculateWidth()
		if(len(files) > i) {
			buttons[i].SetFileInfo(files[i])
		}
		seshLine.AddWidget(buttons[i], 0.125)
	}

	box.SetOrientation(views.Vertical)
	box.AddWidget(seshLine, 0.9)

	app.SetRootWidget(box)
	if e := app.Run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}
}
