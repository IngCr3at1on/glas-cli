package main

import (
	"strings"

	"github.com/IngCr3at1on/glas/ansi"
	tui "github.com/marcusolsson/tui-go"
	"github.com/sirupsen/logrus"
)

const (
	defaultBG = tui.ColorDefault
)

type (
	layout struct {
		log   *logrus.Entry
		ui    tui.UI
		theme *tui.Theme

		inputBox   *tui.Box
		inputEntry *tui.Entry

		outputBox  *tui.Box
		scrollArea *tui.ScrollArea
		scrollBox  *tui.Box

		prevAnsiColors []string
	}
)

func buildLayout(log *logrus.Entry) (l *layout, err error) {
	inputEntry := tui.NewEntry()
	inputEntry.SetFocused(true)
	inputEntry.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputBox := tui.NewHBox(inputEntry)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	outputBox := tui.NewVBox()
	scrollArea := tui.NewScrollArea(outputBox)
	scrollBox := tui.NewVBox(scrollArea)

	l = &layout{
		log:   log,
		theme: tui.NewTheme(),

		inputBox:   inputBox,
		inputEntry: inputEntry,

		outputBox:  outputBox,
		scrollArea: scrollArea,
		scrollBox:  scrollBox,
	}

	l.initTextThemes()
	l.ui, err = tui.New(tui.NewVBox(scrollBox, inputBox))
	if err != nil {
		return nil, err
	}
	l.ui.SetTheme(l.theme)

	return l, nil
}

// Write the contents of p to the tui.Box returning the number of bytes
// written and an error if one exists.
func (l *layout) Write(p []byte) (int, error) {
	// Scroll one line per call.
	l.write(string(p), 1)
	return len(p), nil
}

func (l *layout) write(s string, scroll int) {
	l.outputBox.Append(
		tui.NewHBox(l.ansiText(s)),
	)

	l.scrollArea.Scroll(0, scroll)

	// Run ui.Update asynchronously to avoid blocking.
	go l.ui.Update(func() {
		// Do nothing, we're just forcing a UI update.
	})
}

func (l *layout) initTextThemes() {
	// Hack to be able to setup different colors ahead of time.
	defaultText := tui.NewLabel("")
	defaultText.SetStyleName(ansi.Default)
	l.theme.SetStyle("label.default", tui.Style{Bg: defaultBG, Fg: tui.ColorDefault})

	whiteText := tui.NewLabel("")
	whiteText.SetStyleName(ansi.White)
	l.theme.SetStyle("label.white", tui.Style{Bg: defaultBG, Fg: tui.ColorWhite})

	blackText := tui.NewLabel("")
	blackText.SetStyleName(ansi.Black)
	l.theme.SetStyle("label.black", tui.Style{Bg: defaultBG, Fg: tui.ColorBlack})

	// grey is not supported in tui

	redText := tui.NewLabel("")
	redText.SetStyleName(ansi.Red)
	l.theme.SetStyle("label.red", tui.Style{Bg: defaultBG, Fg: tui.ColorRed})

	greenText := tui.NewLabel("")
	greenText.SetStyleName(ansi.Green)
	l.theme.SetStyle("label.green", tui.Style{Bg: defaultBG, Fg: tui.ColorGreen})

	blueText := tui.NewLabel("")
	blueText.SetStyleName(ansi.Blue)
	l.theme.SetStyle("label.blue", tui.Style{Bg: defaultBG, Fg: tui.ColorBlue})

	yellowText := tui.NewLabel("")
	yellowText.SetStyleName(ansi.Yellow)
	l.theme.SetStyle("label.yellow", tui.Style{Bg: defaultBG, Fg: tui.ColorYellow})

	magentaText := tui.NewLabel("")
	magentaText.SetStyleName(ansi.Magenta)
	l.theme.SetStyle("label.magenta", tui.Style{Bg: defaultBG, Fg: tui.ColorMagenta})

	cyanText := tui.NewLabel("")
	cyanText.SetStyleName(ansi.Cyan)
	l.theme.SetStyle("label.cyan", tui.Style{Bg: defaultBG, Fg: tui.ColorCyan})
}

func (l *layout) ansiText(s string) *tui.Label {
	// FIXME: this doesn't quite work...
	ansiCodes := []string{}
	for k, ac := range ansi.Codes {
		if strings.Contains(s, ac) {
			s = strings.Replace(s, ac, "", -1)
			l.log.Error(k, ac)
			ansiCodes = append(ansiCodes, k)
		}
	}

	var ansiColor string
	if len(ansiCodes) > 0 {
		ansiColor = ansiCodes[0]
		if len(ansiCodes) > 1 {
			if len(l.prevAnsiColors) == 1 {
				if l.prevAnsiColors[0] == ansiCodes[0] {
					ansiColor = ansiCodes[1]
				}
			} else {
				if l.prevAnsiColors[0] != ansiCodes[0] {
					ansiColor = ansiCodes[1]
				}
			}
		}
	}

	label := tui.NewLabel(s)
	if ansiColor != "" {
		label.SetStyleName(ansiColor)

		l.prevAnsiColors = []string{}
		l.prevAnsiColors = append(l.prevAnsiColors, ansiCodes...)
	}

	return label
}
