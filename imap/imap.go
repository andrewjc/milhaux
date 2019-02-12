package imap

import (
	"fmt"
	"github.com/andrewjc/milhaux/common"
	log "github.com/sirupsen/logrus"
	"net"
)

type Imap4Server interface {
	Start(connector common.StorageConnector) error
	handleImapConnection(conn net.Conn)
	closeImapConnection(smtpSession *ImapSession)
}

type Imap4Server_Impl struct {
	config           *common.ApplicationConfig
	commandProcessor ImapCommandProcessor
	listener         net.Listener
	storageConnector common.StorageConnector
}

func NewIMap4Server(config *common.ApplicationConfig) Imap4Server {
	log.Debug("Creating new imap4 server instance...")

	imapSvr := &Imap4Server_Impl{config, NewCommandProcessor(), nil, nil}
	return imapSvr
}

func (s *Imap4Server_Impl) Start(connector common.StorageConnector) error {

	log.Debug("Starting imap server on port ", s.config.GetImap4ServerConfig().Port)

	s.storageConnector = connector

	listenSpec := fmt.Sprintf("%s:%d", s.config.GetImap4ServerConfig().ListenInterface, s.config.GetImap4ServerConfig().Port)
	listener, err := net.Listen("tcp4", listenSpec)

	s.listener = listener

	defer func() {
		s.listener.Close()
		log.Debug("Imap server loop terminated")
	}()

	if err != nil {
		return fmt.Errorf("%s - %s", listenSpec, err.Error())
	}

	log.Info("imap listening on port ", listenSpec)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("error accepting client %s", err.Error())
			continue
		}

		go s.handleImapConnection(conn)
	}

	return nil
}
