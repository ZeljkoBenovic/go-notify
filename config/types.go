package config

import "github.com/hashicorp/go-hclog"

type Config struct {
	MonitoredServices MonitoredServices `yaml:"monitored_services"`
	NotifyService     string            `yaml:"notify_service"`
	Interval          uint64            `yaml:"interval"`
	Timeout           uint64            `yaml:"timeout"`
	ConfigFile        string            `yaml:"config_file,omitempty"`
	Loglevel          string            `yaml:"log_level"`
	LogFileName       string            `yaml:"log_filename"`

	Services NotificationServices `yaml:"notification_services"`

	Logger hclog.Logger `yaml:"logger,omitempty"`
}

type MonitoredServices struct {
	Http []Monitor `yaml:"http"`
}

type Monitor struct {
	Endpoint         string `yaml:"endpoint"`
	ExpectedResponse string `yaml:"expected_response"`
}

type NotificationServices struct {
	Email Email `yaml:"email"`
	Slack Slack `yaml:"slack"`
}

type Email struct {
	To         arrayFlags `yaml:"to"`
	Cc         arrayFlags `yaml:"cc"`
	Bcc        arrayFlags `yaml:"bcc"`
	From       string     `yaml:"from"`
	Subject    string     `yaml:"subject"`
	Body       string     `yaml:"body"`
	AuthUser   string     `yaml:"auth_user"`
	AuthPass   string     `yaml:"auth_pass"`
	SMTPServer string     `yaml:"smtp_server"`
	SMTPPort   uint64     `yaml:"smtp_port"`
	UseAuth    bool       `yaml:"use_auth"`
}

type Slack struct {
	Webhook string `yaml:"webhook"`
}
