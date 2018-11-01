package main

import (
	"os"
	//"fmt"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type SeshButton struct {
	views.Text
	Key rune
	fileInfo os.FileInfo
	RawText string
	textWidth int
}

func (button *SeshButton) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Rune() == button.Key {
			// do callback
			//fmt.Printf(string(button.Key)+" click!\n")
			//fmt.Fprintln(os.Stderr, "three"+string(button.Key))
			button.SetText("CLICK") // temp
			button.RawText = "click"
			button.ReText()
			TestVar := "wtfwtf"
			button.SetText(TestVar)
			//button.Draw()
			//button.Text.PostEventWidgetContent(button)
			//return false
			return true
		}
	}

	return button.Text.HandleEvent(ev)
}

func (button *SeshButton) SetFileInfo(info os.FileInfo) {
	button.fileInfo = info
}

// When SeshButton is resize we need to recalculate textWidth and call ReText()
func (button *SeshButton) Resize() {
	button.CalculateWidth()
	button.ReText()
	button.Text.Resize()
}

// Calculate how wide the text should be based upon the view width
func (button *SeshButton) CalculateWidth() {
	//w, _ := button.view.Size()
	w := 80
	if w < 18 { // cannot do anything reasonable if width < 18
		button.textWidth = 1
	} else {
		button.textWidth = (w-10) / 8
	}
	//fmt.Printf("%d %d\n", w, button.textWidth)
}

// Reset the text by using textWidth and fileInfo
func (button *SeshButton) ReText() {
	button.SetText(button.RawText)
	return

	if button.fileInfo == nil {
		button.RawText = "nil"
		button.SetText("nil")
		return
	}
	filename := button.fileInfo.Name()
	if len(filename) > button.textWidth {
		filename = filename[0:button.textWidth-1]
	}
	button.SetText(button.RawText)
}

//func (button *SeshButton) Draw() {
//	button.Text.Draw()
//}

func NewButton() *SeshButton {
	return &SeshButton{}
}
