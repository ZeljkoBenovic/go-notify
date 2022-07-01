package config

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"os"
)

// errors
var (
	errEndpoint = errors.New("endpoint flag is mandatory")
	errNotify   = errors.New("notify flag is mandatory")
	errResponse = errors.New("response flag is mandatory")
	errToField  = errors.New("email TO field not defined")
)

func cmdGenerateConfigFile() bool {
	if len(os.Args) < 2 {
		return false
	} else {
		cmdArgs := os.Args[1:2]
		return cmdArgs[0] == "config"
	}
}

func (f *Config) checkRequiredData() error {

	if f.MonitoredServices.Http[0].Endpoint == "" {
		return errEndpoint
	}

	if f.NotifyService == "" {
		return errNotify
	}

	if f.MonitoredServices.Http[0].ExpectedResponse == "" {
		return errResponse
	}

	if f.Services.Email.To[0] == "" {
		return errToField
	}

	if f.Services.Email.UseAuth {
		if f.Services.Email.AuthUser == "" {
			f.Services.Email.AuthUser = f.Services.Email.From
		}

		if f.Services.Email.AuthPass == "" {
			//TODO: get smtp pass from env var
			return errors.New("smtp auth password not provided")
		}
	}

	return nil
}

func (f *Config) withDefaults() {
	f.NotifyService = notifyDefault
	f.Interval = intervalDefault
	f.Timeout = timeoutDefault
	f.Loglevel = logLevelDefault

	f.Services.Email.SMTPServer = smtpServerDefault
	f.Services.Email.UseAuth = smtpAuthDefault
	f.Services.Email.SMTPPort = smtpPortDefault
	f.Services.Email.From = emailFromDefault
	f.Services.Email.Subject = emailSubjectDefault

	f.MonitoredServices.Http = []Monitor{
		{
			Endpoint:         "",
			ExpectedResponse: "",
		},
	}
}

func (f *Config) newLogger(name string, logfileLocation string) (hclog.Logger, error) {
	logConfig := &hclog.LoggerOptions{
		Name:  name,
		Level: hclog.LevelFromString(f.Loglevel),
	}

	if f.LogFileName != "" {
		fileWriter, err := os.OpenFile(logfileLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("could not write to log file: %w", err)
		}
		logConfig.Output = fileWriter
	}

	return hclog.New(logConfig), nil
}

func NewConfig() (*Config, error) {
	config := &Config{}
	// get default values
	config.withDefaults()
	// if config arg is set we will generate example config file
	if cmdGenerateConfigFile() {
		if err := config.createConfigFileWithDefaults(); err != nil {
			return nil, fmt.Errorf("could not create default config file %w", err)
		}
		// after we've generated the config file exit the program
		fmt.Println("config file successfully generated")
		os.Exit(0)
	}

	// load the configuration parameters
	if err := config.getConfig(); err != nil {
		return nil, fmt.Errorf("could not load configuration: %w", err)
	}

	// set up logger
	newLogger, err := config.newLogger("go-notify", config.LogFileName)
	if err != nil {
		return nil, fmt.Errorf("could not set up new logger instance: %w", err)
	}

	config.Logger = newLogger

	// check if all required data is present
	if err := config.checkRequiredData(); err != nil {
		config.Logger.Named("command").Error("required parameters not set", "error", err)
		os.Exit(1)
	}

	return config, nil
}
