package main

import (
	"os"

	"github.com/IngCr3at1on/glas"
	tui "github.com/marcusolsson/tui-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type (
	app struct {
		config *Config
		log    *logrus.Entry

		errCh  chan error
		stopCh chan error

		glas   *glas.Glas
		layout *layout

		commandHistory []string
		historyIndex   uint
	}
)

func newApp(path string, config *Config) (*app, error) {
	_app := &app{
		config:         config,
		commandHistory: []string{},
		historyIndex:   0,
	}

	logger := logrus.New()
	logger.SetLevel(config.logLevel)

	file, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "config.LogFile")
	}

	logger.Out = file
	_app.log = logrus.NewEntry(logger)

	_app.layout, err = buildLayout(_app.log)
	if err != nil {
		return nil, err
	}

	_app.errCh = make(chan error)
	_app.stopCh = make(chan error)

	glasConfig := &glas.Config{
		CmdPrefix:     _app.config.CmdPrefix,
		CharacterPath: path,
	}

	_app.glas, err = glas.New(glasConfig, _app.layout, _app.errCh, _app.stopCh, _app.log)
	if err != nil {
		return nil, err
	}
	
	_app.layout.inputEntry.OnSubmit(func(e *tui.Entry) {
		text := e.Text()
		_app.historyIndex = 0
		_app.commandHistory = append(_app.commandHistory, text)

		// Handle any local commands before calling the glas framework.
		ok, err := _app.handleCommand(text)
		if err != nil {
			_app.errCh <- err
			return
		}

		if ok {
			return
		}

		if !_app.config.DisableLocalEcho {
			_app.layout.write(text, 2)
		}
		if err := _app.glas.Send(text); err != nil {
			_app.errCh <- err
			return
		}
		if _app.config.ClearInput {
			e.SetText("")
		}
	})
	_app.layout.ui.SetKeybinding("PgUp", func() {
		_app.layout.scrollArea.Scroll(0, -1)
	})
	_app.layout.ui.SetKeybinding("PgDn", func() {
		_app.layout.scrollArea.Scroll(0, 1)
	})
	_app.layout.ui.SetKeybinding("Up", func() {
		ul := uint(len(_app.commandHistory))
		n := ul - (_app.historyIndex + 1)
		if ul > n && ul > 0 {
			_app.historyIndex++
			_app.layout.inputEntry.SetText(_app.commandHistory[n])
		}
	})
	_app.layout.ui.SetKeybinding("Down", func() {
		// FIXME: this really should clear the input line when you get back to the
		// bottom of the list instead of leaving that last item in the input box.
		if _app.historyIndex > 0 {
			ul := uint(len(_app.commandHistory))
			n := ul - (_app.historyIndex - 1)
			if ul > n && ul > 0 {
				_app.historyIndex--
				_app.layout.inputEntry.SetText(_app.commandHistory[n])
			}
		}
	})

	return _app, nil
}

func (a *app) quit(err error) {
	if err == nil {
		a.stopCh <- errors.New("quit called")
	}

	if err != nil {
		a.log.Infof("quit called with error: %s", err.Error())
	}
}
