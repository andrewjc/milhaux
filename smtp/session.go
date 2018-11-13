package smtp

import (
	"bufio"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

const (
	SESSION_DATA_KEY_CLIENT_ID = "client_id"
	SESSION_DATA_KEY_MAIL_TO   = "to"
	SESSION_DATA_KEY_MAIL_FROM = "from"
	SESSION_DATA_KEY_BUFFER    = "buffer"
)

type SmtpSession struct {
	SmtpServerInstance SmtpServer_impl
	SmtpState          SmtpSessionState
	Writer             *bufio.Writer
	Connection         *net.Conn
	RemoteHostAddr     string
	StateData          map[string]interface{}
	ReceiveDataBuffer  *bytes.Buffer
}

func NewSmtpSession(serverInstance SmtpServer_impl, conn net.Conn) *SmtpSession {
	session := &SmtpSession{
		serverInstance,
		SMTP_SERVER_STATE_ESTABLISH,
		bufio.NewWriter(conn),
		&conn,
		conn.RemoteAddr().String(),
		make(map[string]interface{}),
		bytes.NewBufferString(""),
	}

	return session.beginSession()
}

func (s *SmtpServer_impl) handleSmtpConnection(conn net.Conn) {
	log.Infof("[%s] Accepted smtp connection", conn.RemoteAddr().String())
	smtpSession := NewSmtpSession(*s, conn)
	defer s.closeSmtpConnection(smtpSession)

	rdr := bufio.NewReader(conn)

	for {
		line, err := rdr.ReadString('\n')
		if err != nil {
			log.Errorf("[%s] Error handling smtp connection - ", conn.RemoteAddr().String(), err.Error())
			return
		}

		log.Debugf("[%s] IN: %s", smtpSession.RemoteHostAddr, line)
		log.Debugf("[%s] - raw input: %s", conn.RemoteAddr().String(), string(line))

		commandResponse := s.commandRequestHandler(smtpSession, line)

		switch {
		case commandResponse.action == COMMANDACTION_CONTINUE:
			smtpSession.writeOutputLine(fmt.Sprintf("%d %s", commandResponse.commandStatus, commandResponse.commandResponseText))
			continue
		case commandResponse.action == COMMANDACTION_EXIT:
			smtpSession.writeOutputLine(fmt.Sprintf("%d %s", commandResponse.commandStatus, commandResponse.commandResponseText))
			s.closeSmtpConnection(smtpSession)
			break
		}

	}
}

func (s *SmtpServer_impl) smtpCommandDone(smtpSession *SmtpSession, commandLine string) *CommandResponse {
	return &CommandResponse{COMMANDACTION_EXIT, SMTP_COMMAND_STATUS_SERVICE_CLOSING_CHANNEL, "Bye"}
}

func (s *SmtpServer_impl) closeSmtpConnection(smtpSession *SmtpSession) {
	log.Infof("[%s] Connection closed by remote host (user command)", smtpSession.RemoteHostAddr)
	(*smtpSession.Connection).Close()
}

func (s *SmtpSession) writeOutputLine(outputString string) {
	s.Writer.WriteString(outputString)
	s.Writer.WriteString("\r\n")
	s.Writer.Flush()
	log.Debugf("[%s] OUT: %s", s.RemoteHostAddr, outputString)
}

func (s *SmtpSession) writeOutput(outputString string) {
	s.Writer.WriteString(outputString)
	s.Writer.Flush()
	log.Debugf("[%s] OUT: %s", s.RemoteHostAddr, outputString)
}

func (s *SmtpSession) beginSession() *SmtpSession {
	s.writeOutput(fmt.Sprintf("%d %s ESMTP %s\r\n", SMTP_COMMAND_STATUS_SERVICE_READY, s.SmtpServerInstance.config.GetSmtpServerConfig().Hostname, SMTP_SERVER_BUILD_STRING))
	return s
}
