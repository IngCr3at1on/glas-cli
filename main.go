package main

import (
	"fmt"

	"github.com/IngCr3at1on/glas"
	"github.com/spf13/cobra"
)

const (
	defaultConfigFile = "./config.toml"
	defaultLogFile    = "./glas.log"
	defaultLogLevel   = "error"
)

var (
	configFile string
	logFile    string
	_logLevel  string

	characterFile string

	cmd = &cobra.Command{
		Use:   "glas [address]",
		Short: "A simple MUD Client In Go",

		Run: func(cmd *cobra.Command, args []string) {
			var address string

			if l := len(args); l > 0 && l != 1 {
				fmt.Printf("Invalid number of arguments %d, see --help for more information\n", l)
				return
			} else if l == 1 {
				address = args[0]
			}

			config, err := loadConfig(configFile)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if config == nil {
				if configFile != defaultConfigFile {
					fmt.Printf("config file \"%s\" not found\n", configFile)
				}

				config = &Config{}
			}

			var characterConfig *glas.CharacterConfig
			if characterFile != "" {
				characterConfig, err = loadCharacterConfig(characterFile)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}

			if characterConfig == nil {
				characterConfig = &glas.CharacterConfig{}
			}

			// Prefer the provided address (if there is one).
			if address != "" {
				characterConfig.Address = address
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

			config.logLevel, err = readLogLevel(config.LogLevel)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			_app, err := newApp(config, characterConfig)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			go func() {
				if err := _app.ui.Run(); err != nil {
					_app.errCh <- err
				}
			}()

			// FIXME: from this point forwards all output needs to go to the ui until we exit...

			_app.glas.Start()

		out:
			for {
				select {
				case err = <-_app.stopCh:
					_app.quit(err)
					break out
				case err = <-_app.errCh:
					_app.log.Error(err.Error())
					break out
				}
			}
		},
	}
)

func init() {
	cmd.Flags().StringVar(&configFile, "config", defaultConfigFile, "location for the glas configuration file")
	cmd.Flags().StringVar(&logFile, "logfile", defaultLogFile, "location for log file")
	cmd.Flags().StringVar(&_logLevel, "loglvl", defaultLogLevel, "log level (debug, info, error)")
	cmd.Flags().StringVarP(&characterFile, "character", "c", "", "the character configuration file")
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}
