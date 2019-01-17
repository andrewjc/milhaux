package smtp

import (
	"github.com/andrewjc/milhaux/common"
	"testing"
)

func TestHELOCommand(t *testing.T) {

	s := NewCommandProcessor()

	mockChannel := make(chan *SmtpServerChannelMessage)
	smtpServerInstance := SmtpServer_impl{config: common.CreateDefaultAppConfig(), channel: mockChannel}

	clientConn, _ := getTestConnection()

	mockSession := NewSmtpSession(smtpServerInstance, clientConn, mockChannel)
	commandResponse := s.HandleCommand(mockSession, "helo 127.0.0.1")

	expectSmtpSessionAction(commandResponse, COMMANDACTION_CONTINUE, t)
	expectSmtpStatusCode(commandResponse, SMTP_COMMAND_STATUS_MAIL_ACTION_OK, t)
	expectSmtpStatusText(commandResponse, "127.0.0.1 Hello", t)

}

func TestSubmitsNonHELORequiresHELO(t *testing.T) {

	s := NewCommandProcessor()

	mockChannel := make(chan *SmtpServerChannelMessage)
	smtpServerInstance := SmtpServer_impl{config: common.CreateDefaultAppConfig(), channel: mockChannel}

	clientConn, _ := getTestConnection()

	mockSession := NewSmtpSession(smtpServerInstance, clientConn, mockChannel)

	commandResponse := s.HandleCommand(mockSession, "mail from a@a.com")

	expectSmtpSessionAction(commandResponse, COMMANDACTION_CONTINUE, t)
	expectSmtpStatusCode(commandResponse, SMTP_COMMAND_STATUS_BAD_COMMAND_SEQUENCE, t)
	expectSmtpStatusText(commandResponse, "Send EHLO/HELO first!", t)

}
