package main

import (
	"os"

	"github.com/milhaux/backend"
	"github.com/milhaux/common"
	"github.com/milhaux/imap"
	"github.com/milhaux/smtp"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	appConfig := common.GetAppConfig()

	serverCore := NewApplicationContext(appConfig)
	serverCore.Start()
}

type ApplicationContext struct {
	Config           common.ApplicationConfig
	SmtpServer       smtp.SmtpServer
	ImapServer       imap.Imap4Server
	MailStoreBackend backend.MailStoreBackend
}

func (core *ApplicationContext) Start() {
	log.Info("Milhaux Mail Server.")
	log.Info("Version: ", core.Config.GetApplicationVersion())

	status := core.SmtpServer.Start()
	if status != nil {
		log.Error("Error starting smtp server: ", status.Error())
	}
}

func NewApplicationContext(config *common.ApplicationConfig) *ApplicationContext {
	sc := &ApplicationContext{}
	sc.Config = *config
	sc.SmtpServer = smtp.NewSmtpServer(config)
	sc.ImapServer = imap.NewIMap4Server(config)
	sc.MailStoreBackend = backend.NewMailStoreBackend(config)

	return sc
}
