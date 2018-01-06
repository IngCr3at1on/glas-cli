package main

import (
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	configPath string
	configFile string
	logFile    string
	_logLevel  string

	// TODO: add the ability to load a character file with a specific path.
	//characterFile    string
	clearInput       bool
	disableLocalEcho bool
	cmdPrefix        string

	cmd = &cobra.Command{
		Use:   "glas-cli [connect argument]",
		Short: "A simple MUD Client in Go",
		Long: `A MUD client designed with the terminal in mind. The connect argument can be a
path to a supported character file, the name of a character in the config path,
or an address to connect to.`,

		Run: func(cmd *cobra.Command, args []string) {
			var connectArg string

			if l := len(args); l > 1 {
				fmt.Printf("Invalid number of arguments %d, see --help for more information\n", l)
				return
			} else if l == 1 {
				connectArg = args[0]
			}

			if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
				fmt.Println(err.Error())
				return
			}

			config, err := loadConfig(configFile)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if config.LogFile == "" {
				config.LogFile = logFile
			} else if logFile != defaultLogFile {
				config.LogFile = logFile
			}

			if config.LogLevel == "" {
				config.LogLevel = _logLevel
			} else if _logLevel != defaultLogLevel {
				config.LogLevel = _logLevel
			}

			if clearInput {
				config.ClearInput = clearInput
			}
			if disableLocalEcho {
				config.DisableLocalEcho = disableLocalEcho
			}

			if config.CmdPrefix == "" {
				config.CmdPrefix = cmdPrefix
			} else if cmdPrefix != defaultCmdPrefix {
				config.CmdPrefix = cmdPrefix
			}

			config.logLevel, err = readLogLevel(config.LogLevel)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			_app, err := newApp(configPath, config)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			go func() {
				if err := _app.layout.ui.Run(); err != nil {
					_app.errCh <- err
				}
			}()
			defer func() {
				_app.layout.ui.Quit()
			}()

			_app.glas.Start(connectArg)

		out:
			for {
				select {
				case err = <-_app.stopCh:
					_app.quit(err)
					break out
				case err = <-_app.errCh:
					cause := errors.Cause(err)
					if cause != io.EOF {
						_app.layout.write(err.Error(), 1)
						_app.log.Error(err.Error())
					}
				}
			}
		},
	}
)

func init() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_defaultConfigPath := fmt.Sprintf(defaultConfigPath, usr.HomeDir)
	_defaultConfigFile := fmt.Sprintf(defaultConfigFile, _defaultConfigPath)
	_defaultLogFile := fmt.Sprintf(defaultLogFile, _defaultConfigPath)

	cmd.Flags().StringVarP(&configPath, "path", "p", _defaultConfigPath, "path to search for config files")
	cmd.Flags().StringVar(&configFile, "config", _defaultConfigFile, "location for the glas configuration file")
	cmd.Flags().StringVar(&logFile, "logfile", _defaultLogFile, "location for log file")
	cmd.Flags().StringVar(&_logLevel, "loglvl", defaultLogLevel, "log level (debug, info, error)")
	// TODO: allow loading a specific character by file instead of only the ones loaded by config path.
	//cmd.Flags().StringVarP(&characterFile, "character", "c", "", "the character configuration file")
	cmd.Flags().BoolVar(&clearInput, "clearinput", defaultClearInput, "clear the input bar after hitting enter")
	cmd.Flags().BoolVar(&disableLocalEcho, "localecho", defaultDisableLocalEcho, "whether to display input commands in output")

	if configPath != _defaultConfigPath {
		if configFile == _defaultConfigFile {
			configFile = fmt.Sprintf(defaultConfigFile, configPath)
		}
		if logFile == _defaultLogFile {
			logFile = fmt.Sprintf(defaultLogFile, configPath)
		}
	}
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}
