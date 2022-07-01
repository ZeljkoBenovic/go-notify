package common

import (
	"github.com/ZeljkoBenovic/go-notify/config"
)

// INotifier is an interface that all notification services need to implement
type INotifier interface {
	Send(sendData interface{}) error
	SendMockup() error
	WithConfig(config *config.Config) (INotifier, error)
}

type NotifierFactory func() INotifier

type NotifierType string
