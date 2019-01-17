package smtp

import (
	"fmt"
	"net"

	"github.com/andrewjc/milhaux/common"
	log "github.com/sirupsen/logrus"
)

type SmtpServer interface {
	Start() error
	handleSmtpConnection(conn net.Conn)
	closeSmtpConnection(smtpSession *SmtpSession)

	ObtainListenerChannel() chan *SmtpServerChannelMessage
}

type ChannelMessageType uint8

const (
	SMTP_CHANNEL_MESSAGE_QUEUE_SUBMIT = iota
)

type SmtpServerChannelMessage struct {
	MessageType ChannelMessageType
	Data        *common.MailMessage
}

type SmtpServer_impl struct {
	config           *common.ApplicationConfig
	channel          chan *SmtpServerChannelMessage
	commandProcessor SmtpCommandProcessor

	listener net.Listener
}

func NewSmtpServer(config *common.ApplicationConfig) SmtpServer {
	log.Debug("Creating new smtp server instance...")
	smtpSvr := &SmtpServer_impl{config, make(chan *SmtpServerChannelMessage), NewCommandProcessor(), nil}
	return smtpSvr
}

func (s *SmtpServer_impl) ObtainListenerChannel() chan *SmtpServerChannelMessage {
	return s.channel
}

func (s *SmtpServer_impl) Start() error {

	log.Debug("Starting smtp server on port ", s.config.GetSmtpServerConfig().Port)

	listenSpec := fmt.Sprintf("%s:%d", s.config.GetSmtpServerConfig().ListenInterface, s.config.GetSmtpServerConfig().Port)
	listener, err := net.Listen("tcp4", listenSpec)
	s.listener = listener
	defer func() {
		s.listener.Close()
		log.Debug("SMTP server loop terminated")
	}()

	if err != nil {
		return fmt.Errorf("%s - %s", listenSpec, err.Error())
	}

	log.Info("smtp listening on port ", listenSpec)

	for {
		conn, err := listener.Accept()
		if err != nil {

			log.Errorf("error accepting client %s", err.Error())
			continue
		}

		go s.handleSmtpConnection(conn)
	}

	return nil
}
