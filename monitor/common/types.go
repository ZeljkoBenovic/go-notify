package common

import (
	"github.com/ZeljkoBenovic/go-notify/config"
	"github.com/ZeljkoBenovic/go-notify/notify/common"
)

type IMonitor interface {
	// Run runs the health check against the provided endpoints and sends notifications
	Run() IMonitor
	// RunMock will not send the notifications, used for testing
	RunMock()
	// SetNotifier takes in the notifier interface that monitor will use to send notifications
	SetNotifier(notifier common.INotifier)
}

type MonitorFactory func(config *config.Config) (IMonitor, error)

type MonitorType string
