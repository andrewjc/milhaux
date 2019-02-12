package imap

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/google/uuid"
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

type ImapSession struct {
	imapServerInstance Imap4Server_Impl
	sessionState       ImapSessionState
	writer             *bufio.Writer
	connection         *net.Conn
	remoteHostAddr     string
	stateData          map[string]interface{}
	receiveDataBuffer  *bytes.Buffer
	sessionToken       string
}

func NewImapSession(serverInstance Imap4Server_Impl, conn net.Conn) *ImapSession {
	session := &ImapSession{
		serverInstance,
		IMAP_SESSION_STATE_PREAUTH,
		bufio.NewWriter(conn),
		&conn,
		conn.RemoteAddr().String(),
		make(map[string]interface{}),
		bytes.NewBufferString(""),
		uuid.New().String(),
	}

	return session
}

func (s *Imap4Server_Impl) handleImapConnection(conn net.Conn) {
	log.Infof("[%s] Accepted imap connection", conn.RemoteAddr().String())
	session := NewImapSession(*s, conn)
	defer s.closeImapConnection(session)

	session.beginSession()

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

		log.Debugf("[%s] IN: %s", session.remoteHostAddr, line)

		var commandResponse CommandResponse
		commandResponse = s.commandProcessor.HandleCommand(session, line)

		switch {
		case commandResponse.action == COMMANDACTION_CONTINUE:

			session.writeOutputLines(commandResponse.commandResponseLines)
			continue
		case commandResponse.action == COMMANDACTION_EXIT:
			session.writeOutputLines(commandResponse.commandResponseLines)
			break
		}

	}
}

func (s *Imap4Server_Impl) closeImapConnection(session *ImapSession) {
	log.Infof("[%s] Connection closed by remote host (user command)", session.remoteHostAddr)
	err := (*session.connection).Close()
	if err != nil {
		log.Warnf("[%s] Error occurred while closing imap connection: %s", session.remoteHostAddr, err.Error())
	}
}

func (s *ImapSession) writeOutputLine(outputString string) error {
	return s.writeOutput(fmt.Sprintf("%s\r\n", outputString))
}

func (s *ImapSession) writeOutputLines(lines []string) error {
	for _, line := range lines {
		err := s.writeOutputLine(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ImapSession) writeOutput(outputString string) error {
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

func (s *ImapSession) beginSession() error {
	return s.writeOutput(fmt.Sprintf("* OK Imap server ready for requests from %s (%s)\r\n", s.remoteHostAddr, s.sessionToken))
}
