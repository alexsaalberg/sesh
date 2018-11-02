package main

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

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