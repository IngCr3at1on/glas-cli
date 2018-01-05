package main

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/IngCr3at1on/glas"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type (
	// Config is the app configuration.
	Config struct {
		ClearInput       bool `toml:"clear_input"`
		DisableLocalEcho bool `toml:"disable_local_echo"`
		logLevel         logrus.Level
		LogLevel         string `toml:"log_level"`
		LogFile          string `toml:"log_file"`
	}
)

func loadConfig(file string) (*Config, error) {
	// See if the file exists first.
	if _, err := os.Stat(file); err != nil {
		return nil, nil
	}

	byt, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadFile")
	}

	var config Config
	if _, err := toml.Decode(string(byt), &config); err != nil {
		return nil, errors.Wrap(err, "toml.Decode")
	}

	return &config, nil
}

func loadCharacterConfig(file string) (*glas.CharacterConfig, error) {
	byt, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadFile")
	}

	var config glas.CharacterConfig
	if _, err := toml.Decode(string(byt), &config); err != nil {
		return nil, errors.Wrap(err, "toml.Decode")
	}

	return &config, nil
}

func readLogLevel(lvl string) (logrus.Level, error) {
	switch lvl {
	case "debug":
		return logrus.DebugLevel, nil
	case "info":
		return logrus.InfoLevel, nil
	case "error":
		return logrus.ErrorLevel, nil
	default:
		// We're not actually using fatal anywhere...
		return 0, errors.Errorf("Unrecognized log level: %s (debug, info, error)", lvl)
	}
}
