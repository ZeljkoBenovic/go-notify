package slack

import (
	"errors"
	"fmt"
	"github.com/ZeljkoBenovic/go-http-monitor/config"
	"github.com/ZeljkoBenovic/go-http-monitor/notify/common"
)

type slack struct {
	webhook string
}

func (s slack) SendMockup() error {
	fmt.Println("Sending...")

	return nil
}

func (s slack) WithConfig(config *config.Config) (common.INotifier, error) {
	if config.Services.Slack.Webhook == "" {
		return nil, errors.New("webhook for Slack not defined")
	}
	s.webhook = config.Services.Slack.Webhook

	return s, nil
}

func (s slack) Send(monitor interface{}) error {
	fmt.Println("sending to SLACK webhook ", s.webhook)

	return nil
}

func NotifierFactory() common.INotifier {
	return &slack{}
}
