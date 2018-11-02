package main

import (
	"os"
	"fmt"
	//"reflect"
	//"time"
	//"strconv"
	//"io/ioutil"
	//"io"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

/* Constants & Stuff */

var seshKeys = [8]rune{'a','s','d','f','j','k','l',';'}
var seshColors = [8]tcell.Color{8,9,10,11,12,13,14,15}

/* Sesh Box Stuff */
type SeshBox struct {
	views.BoxLayout
	currentDir string
	App *views.Application

	seshStatus *views.BoxLayout
	seshLine *views.BoxLayout
	seshShell *views.BoxLayout

	seshButtons [8]*SeshButton
	buttonDirsRead int
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
					//buttonFileInfo := b.seshButtons[i].GetFileInfo()
					dirName := b.seshButtons[i].GetFullText()
					if dirName != "" {
						b.navigateToRelativeDir(dirName)
					}
				}
			}
			if ev.Rune() == 'q' {
				b.App.Quit()
				return true
			}
			if ev.Rune() == ' ' {
				fmt.Fprintf(os.Stderr, "space was hit\n")
				b.ShowDirectories()
			}
		}
	}
	return b.BoxLayout.HandleEvent(ev)
}

func exitIfError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func osOpenWrapper(fileName string) *os.File {
	file, err := os.Open(fileName)

	exitIfError(err)

	return file
}

func osReaddirWrapper(dirFile *os.File) []os.FileInfo {
	fileInfos, err := dirFile.Readdir(-1)

	exitIfError(err)

	return fileInfos
}

// populate seshButtons with next 8 directories
func (b *SeshBox) ShowDirectories() {
	numButtonDirs := 7 // number of buttons which get directories (7 cause ';' is gonna manually be '..')

	dirFile, err := os.Open(b.currentDir)
	defer dirFile.Close()
	exitIfError(err)

	fileInfos := osReaddirWrapper(dirFile) // os.FileInfo for each file/dir in b.Currentdir

	var buttonDirs [8]string // the 8(or fewer) directory names (NOT FILES) in b.CurrentDir that will be assigned to sesh buttons

	if b.buttonDirsRead >= len(fileInfos) {
		b.buttonDirsRead = 0
	}

	curDirNum := 0 // how many directories we've found so far
	i := b.buttonDirsRead
	for ; i < len(fileInfos) && curDirNum < numButtonDirs; i++ { // read until we reach end of fileInfos or have found 8 dirs
		if fileInfos[i].Mode().IsDir() {
			buttonDirs[curDirNum] = fileInfos[i].Name()
			curDirNum += 1
		}
	}

	b.buttonDirsRead = i

	for i := 0; i < numButtonDirs; i++ {
		b.seshButtons[i].SetText(buttonDirs[i])
	}

	b.seshButtons[7].SetText("..")

	b.Draw()
}

func (b *SeshBox) Initialize() {
	b.currentDir = "./"

	// // create sesh line
	b.seshLine = views.NewBoxLayout(views.Horizontal)

	var spacer = views.NewText() // add spacer at front edge
	spacer.SetText(" ")
	seshLine.AddWidget(spacer, 0)

	// // create buttons
	for i := 0; i < 8; i++ {
		b.seshButtons[i] = NewButton()
		b.seshButtons[i].SetStyle(tcell.StyleDefault.Background(seshColors[i]))
		b.seshButtons[i].SetAlignment(views.AlignMiddle)
		b.seshButtons[i].Key = seshKeys[i]
		b.seshLine.AddWidget(b.seshButtons[i], 0.1)

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

	// // populate b.SeshButtons with directory names
	b.ShowDirectories()

	// // notify watchers
	b.PostEventWidgetContent(b)
}

func (b *SeshBox) navigateToRelativeDir(dirName string) {
	fullNewPathName := b.currentDir + dirName + string(os.PathSeparator)

	b.currentDir = fullNewPathName
	b.buttonDirsRead = 0
	b.ShowDirectories()
}

func NewSeshBox() *SeshBox {
	return &SeshBox{}
}

/* Sesh Status Stuff */

/* Sesh Button Stuff */
type SeshButton struct {
	views.Text
	Key rune
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
	button.ReText()
}

func (b *SeshButton) SetView(view views.View) {
	b.view = view
	b.Text.SetView(view)
}

func (b *SeshButton) Resize() {
	viewWidth, _ := b.view.Size()
	b.boxWidth = (viewWidth - 10) / 8
	//b.boxWidth = viewWidth
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
// Reset the text by using textWidth and fileInfo
func (button *SeshButton) ReText() {
	filename := button.fullText
	button.fullText = filename
	if len(filename) > button.boxWidth && button.boxWidth > 2 {
		//fmt.Fprintf(os.Stderr,"Clipping %d " +filename+"\n", button.boxWidth)
		filename = filename[0:button.boxWidth-1] + "~"
	}

	button.Text.SetText(filename)
}

func (button *SeshButton) GetFullText() string {
	return button.fullText
}

func NewButton() *SeshButton {
	button := &SeshButton{}
	return button
}
