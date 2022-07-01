package email

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ZeljkoBenovic/go-notify/config"
	monitorHttp "github.com/ZeljkoBenovic/go-notify/monitor/http"
	"github.com/ZeljkoBenovic/go-notify/notify/common"
	"github.com/hashicorp/go-hclog"
	"gopkg.in/gomail.v2"
	"html/template"
	"reflect"
)

type email struct {
	To                                      []string
	cc                                      []string
	bcc                                     []string
	from, subject, body, authUser, authPass string
	smtpAuthEnabled                         bool
	smtpServer                              string
	smtpPort                                uint64
	EndpointsData                           []monitorHttp.HttpEndpoints

	ServiceName []config.Monitor
	smtpMessage *gomail.Message
	logger      hclog.Logger
}

func (e *email) WithConfig(config *config.Config) (common.INotifier, error) {
	if len(config.Services.Email.To) == 0 {
		return nil, errors.New("missing TO address")
	}

	e.body = config.Services.Email.Body
	e.To = config.Services.Email.To
	e.cc = config.Services.Email.Cc
	e.bcc = config.Services.Email.Bcc
	e.from = config.Services.Email.From
	e.subject = config.Services.Email.Subject
	e.authUser = config.Services.Email.AuthUser
	e.authPass = config.Services.Email.AuthPass
	e.smtpAuthEnabled = config.Services.Email.UseAuth
	e.smtpServer = config.Services.Email.SMTPServer
	e.smtpPort = config.Services.Email.SMTPPort
	e.ServiceName = config.MonitoredServices.Http
	e.logger = config.Logger.Named("email")

	e.logger.Debug("email config successfully initialized")
	return e, nil
}

func (e email) Send(monitor interface{}) error {
	//TODO: different reflections for different endpoint structs ( http, json, ... )
	httpType := reflect.ValueOf(monitor)
	httpKind := reflect.Indirect(httpType).Kind()
	if httpKind != reflect.Struct {
		return fmt.Errorf("type Struct must be passed in Send function")
	}

	e.logger.Debug("Send function")

	e.EndpointsData = reflect.Indirect(httpType).FieldByName("Http").Interface().([]monitorHttp.HttpEndpoints)

	if err := e.createMessage(); err != nil {
		return err
	}
	smtp := gomail.NewDialer(e.smtpServer, int(e.smtpPort), e.authUser, e.authPass)

	if err := smtp.DialAndSend(e.smtpMessage); err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}

	e.logger.Info("email notification successfully sent")
	return nil
}

func (e email) SendMockup() error {
	fmt.Println("Sending....")
	return nil
}

func NotifierFactory() common.INotifier {
	return &email{
		smtpMessage: gomail.NewMessage(),
	}
}

func (e *email) createMessage() error {
	e.logger.Debug("createMessage function")

	e.smtpMessage.SetHeader("From", e.from)

	for _, to := range e.To {
		e.smtpMessage.SetHeader("To", to)
	}

	for _, cc := range e.cc {
		e.smtpMessage.SetHeader("Cc", cc)
	}

	for _, bcc := range e.bcc {
		e.smtpMessage.SetHeader("Bcc", bcc)
	}

	e.smtpMessage.SetHeader("Subject", e.subject)

	if e.body == "" {
		if err := e.createHtmlTemplate(); err != nil {
			return fmt.Errorf("could not create email template: %w", err)
		}
	}
	e.smtpMessage.SetBody("text/html", e.body)

	return nil
}

func (e *email) createHtmlTemplate() error {
	buff := new(bytes.Buffer)

	t, err := template.New("default.html").Funcs(template.FuncMap{
		"toString": func(url string) monitorHttp.Url {
			return monitorHttp.Url(url)
		}}).ParseFiles("notify/email/templates/default.html")
	if err != nil {
		return fmt.Errorf("could not parse template %w", err)
	}

	if err := t.Execute(buff, e.EndpointsData); err != nil {
		return fmt.Errorf("could not execute %w", err)
	}

	e.body = buff.String()
	return nil
}
