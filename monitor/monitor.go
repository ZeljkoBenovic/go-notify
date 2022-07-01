package monitor

import (
	"errors"
	"fmt"
	"github.com/ZeljkoBenovic/go-notify/config"
	monitorCommon "github.com/ZeljkoBenovic/go-notify/monitor/common"
	monitorHttp "github.com/ZeljkoBenovic/go-notify/monitor/http"
)

// available notification services
const (
	httpMonitor     monitorCommon.MonitorType = "http"
	telegramMonitor monitorCommon.MonitorType = "telegram"
)

// available services factory
var availableMonitors = map[monitorCommon.MonitorType]monitorCommon.MonitorFactory{
	httpMonitor: monitorHttp.MonitorFactory,
}

func NewMonitor(config *config.Config) (monitorCommon.IMonitor, error) {
	// TODO: Change name for monitor to sender

	// if http monitor
	monitorFactory, ok := availableMonitors["http"]
	if !ok {
		return nil, errors.New("selected monitor type does not exist")
	}

	monitor, err := monitorFactory(config)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate new monitor: %w", err)
	}

	return monitor, nil
}
