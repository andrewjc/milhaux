package core

import (
	"github.com/andrewjc/milhaux/backend"
	"github.com/andrewjc/milhaux/common"
	"github.com/andrewjc/milhaux/imap"
	"github.com/andrewjc/milhaux/smtp"
	log "github.com/sirupsen/logrus"
	"os"
)

type ApplicationContext struct {
	Config           common.ApplicationConfig
	SmtpServer       smtp.SmtpServer
	ImapServer       imap.Imap4Server
	MailStoreBackend *backend.MailStoreBackend
}

func (core *ApplicationContext) Start() {
	log.Info("Milhaux Mail Server.")
	log.Info("Version: ", core.Config.GetApplicationVersion())

	smtpServerChannel := core.SmtpServer.ObtainListenerChannel()

	go core.initBackend(smtpServerChannel)

	go core.initSmtpServer()

	core.beginMainMessageLoop(common.GetMainMessageLoop())
}

func (core *ApplicationContext) beginMainMessageLoop(messageChannel chan common.MainEventMessage) {
	for {
		select {
		case eventMessage := <-messageChannel:
			if eventMessage.MessageType == common.SHUTDOWN {
				os.Exit(0)
			}
			if eventMessage.MessageType == common.PING {
				log.Info("Ping...")
			}
		}
	}
}

func (core *ApplicationContext) initSmtpServer() {
	if status := core.SmtpServer.Start(); status != nil {
		log.Error("Error starting smtp server: ", status.Error())
		return
	}
}

func (core *ApplicationContext) initBackend(messageChannel chan smtp.SmtpServerChannelMessage) {

	core.MailStoreBackend.InitSmtpMessageChannelListener(messageChannel)

	if status := core.MailStoreBackend.Start(); status != nil {
		log.Error("Error starting mailer backend:", status.Error())
		return
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
