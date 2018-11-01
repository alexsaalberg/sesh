package main

import (
	"os"
	"fmt"
	//"reflect"
	//"time"
	//"strconv"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

/* Sesh Box Stuff */
type SeshBox struct {
	views.BoxLayout
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

// overloaded functions
func (button *SeshButton) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Rune() == button.Key {
			// do callback
			//button.SetText("CLICK")

			//button.Text.PostEventWidgetContent(button)
			//return true
		}
	}

	return button.Text.HandleEvent(ev)
}

func (button *SeshButton) SetText(text string) {
	button.fullText = text
	//if len(text) > button.boxWidth && button.boxWidth > 2{
	//	text = text[0:button.boxWidth - 2] + "~"
	//}
	button.Text.SetText(text)
}

func (b *SeshButton) SetView(view views.View) {
	b.view = view
	b.Text.SetView(view)
}

func (b *SeshButton) Resize() {
	viewWidth, _ := b.view.Size()
	b.boxWidth = (viewWidth - 10) / 8
}

func (b *SeshButton) Size() (int, int) {
	b.Resize()
	b.ReText()
	return b.boxWidth, 1
}

func (b *SeshButton) Draw() {
	b.Resize()
	b.ReText()
	b.Text.Draw()
}

// new functions
func (button *SeshButton) SetFileInfo(info os.FileInfo) {
	button.fileInfo = info
	button.ReText()
}

// Reset the text by using textWidth and fileInfo
func (button *SeshButton) ReText() {
	//fmt.Fprintf(os.Stderr, "TextWidth: "+strconv.Itoa(button.boxWidth)+"\n")
	if button.fileInfo == nil || button.boxWidth <= 0 {
		button.SetText("nil")
		return
	}

	filename := button.fileInfo.Name()
	button.fullText = filename
	if len(filename) > button.boxWidth && button.boxWidth > 2 {
		fmt.Fprintf(os.Stderr,"Clipping %d " +filename+"\n", button.boxWidth)
		filename = filename[0:button.boxWidth-2] + "~"
	}

	button.SetText(filename)
}

func NewButton() *SeshButton {
	button := &SeshButton{}
	return button
}
