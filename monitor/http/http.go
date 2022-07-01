package http

import (
	"fmt"
	"github.com/ZeljkoBenovic/go-notify/config"
	"github.com/ZeljkoBenovic/go-notify/monitor/common"
	notifyCommon "github.com/ZeljkoBenovic/go-notify/notify/common"
	"github.com/hashicorp/go-hclog"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Url string

type HttpMonitor struct {
	Http []HttpEndpoints

	Timeout uint64
	Logger  hclog.Logger

	Sender notifyCommon.INotifier
}

type HttpEndpoints struct {
	Url          string
	Error        error
	SearchString map[Url]string
	Result       map[Url]string
	Healthy      map[Url]bool
}

//MonitorFactory is the factory method for http monitor
func MonitorFactory(config *config.Config) (common.IMonitor, error) {
	mon := &HttpMonitor{}

	mon.Timeout = config.Timeout
	mon.Logger = config.Logger

	for _, srvc := range config.MonitoredServices.Http {
		mon.Http = append(
			mon.Http,
			HttpEndpoints{
				Url:          srvc.Endpoint,
				Result:       map[Url]string{Url(srvc.Endpoint): ""},
				SearchString: map[Url]string{Url(srvc.Endpoint): srvc.ExpectedResponse},
				Healthy:      map[Url]bool{Url(srvc.Endpoint): false},
			})
	}
	return mon, nil
}

// SetNotifier sets the sender notifier interface
func (m *HttpMonitor) SetNotifier(notifier notifyCommon.INotifier) {
	m.Sender = notifier
}

// Run runs the health check and sends notifications
func (m *HttpMonitor) Run() common.IMonitor {
	wg := sync.WaitGroup{}
	mux := sync.Mutex{}

	// fetch endpoints
	for _, httpEndpoint := range m.Http {
		wg.Add(1)

		// query endpoints in parallel
		go func(httpEndpoint HttpEndpoints) {
			defer wg.Done()

			client := http.Client{Timeout: time.Duration(m.Timeout) * time.Second}

			resp, err := client.Get(httpEndpoint.Url)
			if err != nil {
				m.Logger.Debug("could not send request to", "url", httpEndpoint.Url)
				mux.Lock()
				httpEndpoint.Error = fmt.Errorf("could not send GET request err=%w", err)
				mux.Unlock()
			}
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				m.Logger.Debug("response body", "body", string(body))
				mux.Lock()
				httpEndpoint.Error = fmt.Errorf("could not read the responce body err=%w", err)
				mux.Unlock()
			}

			mux.Lock()
			httpEndpoint.Result[Url(httpEndpoint.Url)] = string(body)
			mux.Unlock()

			m.Logger.Info("successfully queried defined url", "url", httpEndpoint.Url)
		}(httpEndpoint)
	}

	wg.Wait()

	m.checkHealth().writeLogs().sendNotifications()
	return m
}

// checkHealth checks the data received and creates a map with bool health values
func (m *HttpMonitor) checkHealth() *HttpMonitor {
	// check if the strings are found in results
	for _, e := range m.Http {
		e.Healthy[Url(e.Url)] = strings.Contains(e.Result[Url(e.Url)], e.SearchString[Url(e.Url)])
		switch e.Healthy[Url(e.Url)] {
		case true:
			m.Logger.Info("service health", "url", e.Url, "status", "HEALTHY")
		case false:
			m.Logger.Info("service health", "url", e.Url, "status", "NOT-HEALTHY")
		}
	}

	return m
}

// writeLogs writes the logs in the console
func (m HttpMonitor) writeLogs() *HttpMonitor {
	for _, mon := range m.Http {
		if !mon.Healthy[Url(mon.Url)] {
			m.Logger.Warn("Service health entered ALARM state", "url", mon.Url)
		} else {
			m.Logger.Info("Service HEALTHY", "url", mon.Url)
		}
	}

	return &m
}

// RunMock doesn't send any notifications
func (m HttpMonitor) RunMock() {
	for _, mon := range m.Http {

		if !mon.Healthy[Url(mon.Url)] {

			m.Logger.Info("Sending mock notifications...", "url", mon.Url)

			if err := m.Sender.SendMockup(); err != nil {
				m.Logger.Error("Could not send mock notifications", "url", mon.Url, "error", err.Error())
			}

			// we break from send notification on first failed as one notification is enough
			break
		}
	}
}

func (m *HttpMonitor) sendNotifications() {
	for _, mon := range m.Http {

		if !mon.Healthy[Url(mon.Url)] {

			m.Logger.Info("Sending notifications...", "url", mon.Url)

			if err := m.Sender.Send(m); err != nil {
				m.Logger.Error("Could not send notifications", "url", mon.Url, "error", err.Error())
			}
			// we break from send notification on first failed as one notification is enough
			break
		}
	}
}
