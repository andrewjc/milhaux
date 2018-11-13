package smtp

import (
	"fmt"
	"net"

	"github.com/milhaux/common"
	log "github.com/sirupsen/logrus"
)

type SmtpServer interface {
	Start() error
	handleSmtpConnection(conn net.Conn)
	closeSmtpConnection(smtpSession *SmtpSession)
}

type SmtpServer_impl struct {
	config *common.ApplicationConfig

	listener net.Listener
}

func NewSmtpServer(config *common.ApplicationConfig) SmtpServer {

	log.Debug("Creating new smtp server instance...")
	smtpSvr := &SmtpServer_impl{config, nil}
	return smtpSvr
}

func (s *SmtpServer_impl) Start() error {

	log.Debug("Starting smtp server on port ", s.config.GetSmtpServerConfig().Port)

	listenSpec := fmt.Sprintf("%s:%d", s.config.GetSmtpServerConfig().ListenInterface, s.config.GetSmtpServerConfig().Port)
	listener, err := net.Listen("tcp4", listenSpec)

	s.listener = listener

	if err != nil {
		return fmt.Errorf("%s - %s", listenSpec, err.Error())
	}

	log.Info("smtp listening on port ", listenSpec)

	defer s.listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {

			log.Error("error accepting client %s", err.Error())
			continue
		}

		go s.handleSmtpConnection(conn)
	}

	return nil
}
func (s *SmtpServer_impl) onSubmitMail(message *MailMessage) *MailMessage {
	message.queueId = "123123123"
	return message
}
