package main

import (
	"os"
	"sync"

	"github.com/IngCr3at1on/glas"
	tui "github.com/marcusolsson/tui-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var quitOnce sync.Once

type (
	app struct {
		config          *Config
		characterConfig *glas.CharacterConfig
		log             *logrus.Entry

		errCh  chan error
		stopCh chan error

		glas   *glas.Glas
		layout *layout

		ui tui.UI
	}
)

func newApp(config *Config, characterConfig *glas.CharacterConfig) (*app, error) {
	_app := &app{}

	logger := logrus.New()
	logger.SetLevel(config.logLevel)

	file, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "config.LogFile")
	}

	logger.Out = file
	_app.log = logrus.NewEntry(logger)

	_app.layout = buildLayout(_app.log)

	_app.errCh = make(chan error)
	_app.stopCh = make(chan error)

	_app.glas, err = glas.New(characterConfig, _app.layout.output, _app.errCh, _app.stopCh, _app.log)
	if err != nil {
		return nil, err
	}

	_app.layout.input.entry.OnSubmit(func(e *tui.Entry) {
		if err := _app.glas.Send(e.Text()); err != nil {
			_app.errCh <- err
		}

		// TODO: control this in settings.
		_app.layout.output.write(e.Text(), 2)

		// FIXME: clear input box (having issues doing this)...
	})

	_app.ui = tui.New(tui.NewVBox(_app.layout.output.scrollBox, _app.layout.input.box))

	_app.ui.SetKeybinding("Esc", func() {
		_app.quit(nil)
	})
	_app.ui.SetKeybinding("Up", func() {
		_app.layout.output.scrollArea.Scroll(0, -1)
	})
	_app.ui.SetKeybinding("Down", func() {
		_app.layout.output.scrollArea.Scroll(0, 1)
	})

	return _app, nil
}

func (a *app) quit(err error) {
	quitOnce.Do(func() {
		a.ui.Quit()
	})

	if err == nil {
		a.stopCh <- errors.New("quit called")
	}

	if err != nil {
		a.log.Infof("quit called with error: %s", err.Error())
	}
}
