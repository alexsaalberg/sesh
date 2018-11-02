package main

import (
	"os"
	"fmt"
	//"reflect"
	//"time"
	//"strconv"
	//"io/ioutil"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

/* Constants & Stuff */

var seshKeys = [8]rune{'a','s','d','f','j','k','l',';'}
var seshColors = [8]tcell.Color{8,9,10,11,12,13,14,15}

/* Sesh Box Stuff */
type SeshBox struct {
	views.BoxLayout
	currentDir *os.File
	App *views.Application

	seshStatus *views.BoxLayout
	seshLine *views.BoxLayout
	seshShell *views.BoxLayout

	buttons [8]*SeshButton
}

func (b *SeshBox) HandleEvent (ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape, tcell.KeyEnter:
			b.App.Quit()
			return true
		case tcell.KeyRune:
			for i, key := range seshKeys {
				if ev.Rune() == key {
					buttonFileInfo := b.buttons[i].GetFileInfo()

					if buttonFileInfo != nil {
						b.navigateToDir(buttonFileInfo)
					}
				}
			}
			if ev.Rune() == 'q' {
				b.App.Quit()
				return true
			}
		}
	}
	return b.BoxLayout.HandleEvent(ev)
}

func (b *SeshBox) Initialize() {
	// // get current dir as os.File
    dirname := "." + string(os.PathSeparator)
    d, err := os.Open(dirname)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
	b.currentDir = d

	fmt.Fprintf(os.Stderr, d.Name())

    fi, err := d.Readdir(-1)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

	// // get array of dirs
	var subDirs [8]os.FileInfo
	dirNum := 0
    for _, fi := range fi {
        if dirNum < 8 && fi.Mode().IsDir() {
			subDirs[dirNum] = fi
            //fmt.Fprintln(os.Stderr, fi.Name(), fi.Size(), "bytes")
        }
    }

	// // create sesh line
	b.seshLine = views.NewBoxLayout(views.Horizontal)

	var spacer = views.NewText() // add spacer at front edge
	spacer.SetText(" ")
	seshLine.AddWidget(spacer, 0)

	// // create buttons
	for i := 0; i < 8; i++ {
		b.buttons[i] = NewButton()
		b.buttons[i].SetStyle(tcell.StyleDefault.Background(tcell.ColorGrey))
		b.buttons[i].SetAlignment(views.AlignMiddle)
		b.buttons[i].Key = seshKeys[i]
		b.buttons[i].SetFileInfo(subDirs[i]) // assign dir to button
		b.seshLine.AddWidget(b.buttons[i], 0.1)

		// add spacer between buttons
		spacer = views.NewText()
		spacer.SetText(" ")
		if i == 3 {
			spacer.SetText("  ")
		}
		b.seshLine.AddWidget(spacer, 0)
	}

	// // create textbox for ls output
	b.seshStatus = views.NewBoxLayout(views.Horizontal)

	// // create shell
	b.seshShell = views.NewBoxLayout(views.Horizontal)

	// // add all the widgets
	b.AddWidget(b.seshStatus, 0.5)
	b.AddWidget(b.seshLine, 0)
	b.AddWidget(b.seshShell, 0.5)

	// render ls output

	// // notify watchers
	b.PostEventWidgetContent(b)
}

func (b *SeshBox) navigateToDir(newDir os.FileInfo) {
	newDirName := b.currentDir.Name() + newDir.Name() + string(os.PathSeparator) 
	//newDirName := newDir.Name()

    dir, err := os.Open(newDirName)
    if err != nil {
		fmt.Fprintf(os.Stderr, "Error in navigateToDir; oldDir: %q, newDir: %q, fullNewPath: %q", b.currentDir.Name(), newDir.Name(), newDirName)
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

	b.currentDir = dir

    fi, err := dir.Readdir(-1)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

	var subDirs [8]os.FileInfo
	dirNum := 0
    for _, fi := range fi {
        if dirNum < 8 && fi.Mode().IsDir() {
			subDirs[dirNum] = fi
			dirNum += 1
            //fmt.Fprintln(os.Stderr, fi.Name(), fi.Size(), "bytes")
        }
    }

	// make ; go up one dir
	subDirs[7], err = os.Stat("..") 
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for i, subDir := range subDirs {
		b.buttons[i].SetFileInfo(subDir)
	}

	b.Draw()
	b.PostEventWidgetContent(b)
}

func NewSeshBox() *SeshBox {
	return &SeshBox{}
}

/* Sesh Status Stuff */

/* Sesh Button Stuff */
type SeshButton struct {
	views.Text
	Key rune
	fileInfo os.FileInfo
	boxWidth int
	fullText string // unclipped text
	view views.View
}


/// Overloaded functions

func (button *SeshButton) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Rune() == button.Key {
			//fmt.Fprintf(os.Stderr, ".%c clicked.", button.Key)
		}
	}

	return button.Text.HandleEvent(ev)
}

func (button *SeshButton) SetText(text string) {
	button.fullText = text
	button.Text.SetText(text)
}

func (b *SeshButton) SetView(view views.View) {
	b.view = view
	b.Text.SetView(view)
}

func (b *SeshButton) Resize() {
	viewWidth, _ := b.view.Size()
	b.boxWidth = (viewWidth - 10) / 8
	//fmt.Fprintf(os.Stderr, "Resize(): %d %d\n", b.boxWidth, viewWidth)
	b.boxWidth = viewWidth
	b.Text.Resize()
}

func (b *SeshButton) Size() (int, int) {
	b.Resize()
	b.ReText()
	return b.boxWidth, 3
}

func (b *SeshButton) Draw() {
	b.Resize()
	b.ReText()
	b.Text.Draw()
}


/// New functions

func (button *SeshButton) SetFileInfo(info os.FileInfo) {
	button.fileInfo = info
	button.ReText()
}

func (button *SeshButton) GetFileInfo() (info os.FileInfo) {
	return button.fileInfo
}

// Reset the text by using textWidth and fileInfo
func (button *SeshButton) ReText() {
	//fmt.Fprintf(os.Stderr, "TextWidth: "+strconv.Itoa(button.boxWidth)+"\n")
	if button.fileInfo == nil {
		button.SetText("")
		return
	}

	filename := button.fileInfo.Name()
	button.fullText = filename
	if len(filename) > button.boxWidth && button.boxWidth > 2 {
		//fmt.Fprintf(os.Stderr,"Clipping %d " +filename+"\n", button.boxWidth)
		filename = filename[0:button.boxWidth-1] + "~"
	}

	button.SetText(filename)
}

func NewButton() *SeshButton {
	button := &SeshButton{}
	return button
}
