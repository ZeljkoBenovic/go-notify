package config

import (
	"flag"
	"fmt"
)

// default config values
const (
	notifyDefault   string = "email"
	intervalDefault uint64 = 300
	timeoutDefault  uint64 = 60

	smtpAuthDefault     bool   = false
	smtpServerDefault   string = "localhost"
	smtpPortDefault     uint64 = 25
	emailFromDefault    string = "gonotify@service.check"
	emailSubjectDefault string = "[GONOTIFY] SERVICE ENTERED AN ALARM STATE"
	logLevelDefault     string = "INFO"
)

type arrayFlags []string

func (a arrayFlags) String() string {
	return "array flag"
}

func (a *arrayFlags) Set(s string) error {
	*a = append(*a, s)
	return nil
}

func (f *Config) getConfig() error {
	flag.StringVar(&f.ConfigFile, "config", "", "Config file location ( .yaml, .yml )")
	flag.StringVar(&f.NotifyService, "notify", notifyDefault, "Service used for notification to notify (email, slack)")
	//flag.StringVar(&f.MonitoredServices.Http[0].Endpoint, "endpoint", "", "Endpoint to monitor")
	//flag.StringVar(&f.Response, "resp-str", "", "Expected string in response")
	flag.Uint64Var(&f.Interval, "interval", intervalDefault, "Interval in seconds to query the endpoint")
	flag.Uint64Var(&f.Timeout, "timeout", timeoutDefault, "Timeout in seconds to consider an endpoint unresponsive")
	flag.StringVar(&f.Loglevel, "log-level", logLevelDefault, "Log level output (INFO, DEBUG)")
	flag.StringVar(&f.LogFileName, "log-file", "", "Log file name to output all logs")

	flag.BoolVar(&f.Services.Email.UseAuth, "smtp-auth", smtpAuthDefault, "Set to true if your SMTP server requires SMTP authentication")
	flag.StringVar(&f.Services.Email.SMTPServer, "smtp-server", smtpServerDefault, "SMTP server and port that will be used to send email")
	flag.Uint64Var(&f.Services.Email.SMTPPort, "smtp-port", smtpPortDefault, "SMTP server port")
	flag.Var(&f.Services.Email.To, "email-to", "Email addresses to send the notification")
	flag.Var(&f.Services.Email.Cc, "email-cc", "Email addresses for cc field")
	flag.Var(&f.Services.Email.Bcc, "email-bcc", "Email addresses bcc field")
	flag.StringVar(&f.Services.Email.From, "email-from", emailFromDefault, "Email address from which to send the notification")
	flag.StringVar(&f.Services.Email.Subject, "email-subject", emailSubjectDefault, "Subject of the notification email")
	flag.StringVar(&f.Services.Email.Body, "email-body", "", "Body of the notification email")
	flag.StringVar(&f.Services.Email.AuthUser, "smtp-user", "", "SMTP user used for authentication")
	flag.StringVar(&f.Services.Email.AuthPass, "smtp-pass", "", "SMTP pass used for authentication")
	flag.Parse()

	if f.ConfigFile != "" {
		if err := f.loadFromConfigFile(); err != nil {
			return fmt.Errorf("could not load config from file: %w", err)
		}
	}

	return nil
}
