package main

import (
	tui "github.com/marcusolsson/tui-go"
	"github.com/sirupsen/logrus"
)

type (
	layout struct {
		input  *input
		output *output
		ui     tui.UI
	}

	input struct {
		box   *tui.Box
		entry *tui.Entry
	}

	output struct {
		log        *logrus.Entry
		box        *tui.Box
		scrollArea *tui.ScrollArea
		scrollBox  *tui.Box
	}
)

func buildLayout(log *logrus.Entry) *layout {
	output := &output{
		log: log,
		box: tui.NewVBox(),
	}
	output.scrollArea = tui.NewScrollArea(output.box)
	output.scrollBox = tui.NewVBox(output.scrollArea)

	input := &input{
		entry: tui.NewEntry(),
	}
	input.entry.SetFocused(true)
	input.entry.SetSizePolicy(tui.Expanding, tui.Maximum)
	input.box = tui.NewHBox(input.entry)
	input.box.SetSizePolicy(tui.Expanding, tui.Maximum)

	return &layout{
		output: output,
		input:  input,
		ui:     tui.New(tui.NewVBox(output.scrollBox, input.box)),
	}
}

// Write the contents of p to the tui.Box returning the number of bytes
// written and an error if one exists.
func (l *layout) Write(p []byte) (int, error) {
	// Scroll one line per call.
	l.write(string(p), 1)
	return len(p), nil
}

func (l *layout) write(s string, scroll int) {
	l.output.box.Append(tui.NewHBox(
		tui.NewLabel(s),
	))

	l.output.scrollArea.Scroll(0, scroll)

	l.ui.Update(func() {
		// Do nothing, we're just forcing a UI update.
	})
}
