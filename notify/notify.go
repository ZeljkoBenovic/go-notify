package notify

import (
	"errors"
	"fmt"
	"github.com/ZeljkoBenovic/go-notify/config"
	"github.com/ZeljkoBenovic/go-notify/notify/common"
	"github.com/ZeljkoBenovic/go-notify/notify/email"
	"github.com/ZeljkoBenovic/go-notify/notify/slack"
)

// available notifier service names
const (
	slackType common.NotifierType = "slack"
	emailType common.NotifierType = "email"
)

// availableNotifiers creates a map of all available NotifierFactories
var availableNotifiers = map[common.NotifierType]common.NotifierFactory{
	emailType: email.NotifierFactory,
	slackType: slack.NotifierFactory,
}

// NewNotifier returns an instance of the notifier service
func NewNotifier(config *config.Config) (common.INotifier, error) {
	notifierFactory, ok := availableNotifiers[common.NotifierType(config.NotifyService)]
	if !ok {
		return nil, errors.New("selected notifier not available")
	}

	notifierService, err := notifierFactory().WithConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create notifier instance: %w", err)
	}

	return notifierService, nil
}
