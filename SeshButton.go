package main

import (
	"os"
	
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type SeshButton struct {
	Key rune
	fileInfo os.FileInfo
	textWidth int
	views.Text
}

func (button *SeshButton) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Rune() == button.Key {
			// do callback
			//fmt.Printf(string(button.Key)+" click!\n")
			button.SetText("CLICK") // temp
			button.Text.PostEventWidgetContent(button)
			return true
		}
	}
	return button.Text.HandleEvent(ev)
}

func (button *SeshButton) SetFileInfo(info os.FileInfo) {
	button.fileInfo = info
	button.ReText()
}

// When SeshButton is resize we need to recalculate textWidth and call ReText()
func (button *SeshButton) Resize() {
	button.CalculateWidth()
	button.ReText()
	button.Text.Resize()
}

// Calculate how wide the text should be based upon the view width
func (button *SeshButton) CalculateWidth() {
	//w, h := button.Text.view.Size()
	w := 80
	if w < 18 { // cannot do anything reasonable if width < 18
		button.textWidth = 1
	} else {
		button.textWidth = (w-10) / 8
	}
}

// Reset the text by using textWidth and fileInfo
func (button *SeshButton) ReText() {
	if button.fileInfo == nil {
		button.SetText("nil")
		return
	}
	filename := button.fileInfo.Name()
	if len(filename) > button.textWidth {
		filename = filename[0:button.textWidth-1]
	}
	button.SetText(filename)
}

func NewButton() *SeshButton {
	return &SeshButton{}
}
