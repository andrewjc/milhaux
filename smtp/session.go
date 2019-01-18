package smtp

import (
	"bufio"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

const (
	SESSION_DATA_KEY_CLIENT_ID = "client_id"
	SESSION_DATA_KEY_MAIL_TO   = "to"
	SESSION_DATA_KEY_MAIL_FROM = "from"
	SESSION_DATA_KEY_BUFFER    = "buffer"
)

type SmtpSession struct {
	smtpServerInstance SmtpServer_impl
	smtpState          SmtpSessionState
	writer             *bufio.Writer
	connection         *net.Conn
	remoteHostAddr     string
	stateData          map[string]interface{}
	receiveDataBuffer  *bytes.Buffer
	smtpMessageChannel chan SmtpServerChannelMessage
}

func NewSmtpSession(serverInstance SmtpServer_impl, conn net.Conn, messageChannel chan SmtpServerChannelMessage) *SmtpSession {
	session := &SmtpSession{
		serverInstance,
		SMTP_SESSION_STATE_PREAUTH,
		bufio.NewWriter(conn),
		&conn,
		conn.RemoteAddr().String(),
		make(map[string]interface{}),
		bytes.NewBufferString(""),
		messageChannel,
	}

	return session
}

func (s *SmtpServer_impl) handleSmtpConnection(conn net.Conn) {
	log.Infof("[%s] Accepted smtp connection", conn.RemoteAddr().String())
	smtpSession := NewSmtpSession(*s, conn, s.ObtainListenerChannel())
	defer s.closeSmtpConnection(smtpSession)

	smtpSession.beginSession()

	rdr := bufio.NewReader(conn)

	for {
		line, err := rdr.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Errorf("[%s] Error handling smtp connection - %s", conn.RemoteAddr().String(), err.Error())
			break
		}

		log.Debugf("[%s] IN: %s", smtpSession.remoteHostAddr, line)

		var commandResponse CommandResponse
		commandResponse = s.commandProcessor.HandleCommand(smtpSession, line)

		switch {
		case commandResponse.action == COMMANDACTION_CONTINUE:
			smtpSession.writeOutputLine(fmt.Sprintf("%d %s", commandResponse.commandStatus, commandResponse.commandResponseText))
			continue
		case commandResponse.action == COMMANDACTION_EXIT:
			smtpSession.writeOutputLine(fmt.Sprintf("%d %s", commandResponse.commandStatus, commandResponse.commandResponseText))
			break
		}

	}
}

func (s *SmtpServer_impl) closeSmtpConnection(smtpSession *SmtpSession) {
	log.Infof("[%s] Connection closed by remote host (user command)", smtpSession.remoteHostAddr)
	err := (*smtpSession.connection).Close()
	if err != nil {
		log.Warnf("[%s] Error occurred while closing smtp connection: %s", smtpSession.remoteHostAddr, err.Error())
	}
}

func (s *SmtpSession) writeOutputLine(outputString string) error {
	return s.writeOutput(fmt.Sprintf("%s\n", outputString))
}

func (s *SmtpSession) writeOutput(outputString string) error {
	_, err := s.writer.WriteString(outputString)
	if err != nil {
		log.Debugf("[%s] writeOutput error: %s", s.remoteHostAddr, err.Error())
		return err
	}

	err = s.writer.Flush()
	if err != nil {
		log.Debugf("[%s] writeOutput flush error: %s", s.remoteHostAddr, err.Error())
		return err
	}

	log.Debugf("[%s] OUT: %s", s.remoteHostAddr, outputString)

	return nil
}

func (s *SmtpSession) beginSession() error {
	return s.writeOutput(fmt.Sprintf("%d %s ESMTP %s\r\n", SMTP_COMMAND_STATUS_SERVICE_READY, s.smtpServerInstance.config.GetSmtpServerConfig().Hostname, SMTP_SERVER_BUILD_STRING))
}
