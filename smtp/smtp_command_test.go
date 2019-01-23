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

	commandResponse := s.HandleCommand(mockSession, "mail from: a@a.com")

	expectSmtpSessionAction(commandResponse, COMMANDACTION_CONTINUE, t)
	expectSmtpStatusCode(commandResponse, SMTP_COMMAND_STATUS_BAD_COMMAND_SEQUENCE, t)
	expectSmtpStatusText(commandResponse, "Send EHLO/HELO first!", t)

}

func TestAllowMultipleMessagePerSession(t *testing.T) {

	s := NewCommandProcessor()

	mockChannel := make(chan *SmtpServerChannelMessage)
	smtpServerInstance := SmtpServer_impl{config: common.CreateDefaultAppConfig(), channel: mockChannel}
	smtpServerInstance.config.GetSmtpServerConfig().SMTP_OPTION_SINGLE_MESSAGE_PER_SESSION = true

	clientConn, _ := getTestConnection()

	mockSession := NewSmtpSession(smtpServerInstance, clientConn, mockChannel)

	mockSession.beginSession()
	s.HandleCommand(mockSession, "helo 127.0.0.1")

	// fast forward the session
	mockSession.smtpState = SMTP_SESSION_STATE_SUBMIT

	commandResponse := s.HandleCommand(mockSession, "mail from: a@a.com")

	// check that the session has reverted it's state to the start
	if mockSession.smtpState != SMTP_SESSION_STATE_MAIL {
		t.Errorf("SMTP state expected SMTP_SESSION_STATE_MAIL")
	}

	expectSmtpSessionAction(commandResponse, COMMANDACTION_CONTINUE, t)
	expectSmtpStatusCode(commandResponse, SMTP_COMMAND_STATUS_MAIL_ACTION_OK, t)
	expectSmtpStatusText(commandResponse, "OK", t)
}

func TestDisAllowMultipleMessagePerSession(t *testing.T) {

	s := NewCommandProcessor()

	mockChannel := make(chan *SmtpServerChannelMessage)
	smtpServerInstance := SmtpServer_impl{config: common.CreateDefaultAppConfig(), channel: mockChannel}
	smtpServerInstance.config.GetSmtpServerConfig().SMTP_OPTION_SINGLE_MESSAGE_PER_SESSION = false

	clientConn, _ := getTestConnection()

	mockSession := NewSmtpSession(smtpServerInstance, clientConn, mockChannel)

	mockSession.beginSession()
	s.HandleCommand(mockSession, "helo 127.0.0.1")

	// fast forward the session
	mockSession.smtpState = SMTP_SESSION_STATE_SUBMIT

	commandResponse := s.HandleCommand(mockSession, "mail from a@a.com")

	expectSmtpSessionAction(commandResponse, COMMANDACTION_EXIT, t)
	expectSmtpStatusCode(commandResponse, SMTP_COMMAND_STATUS_SERVICE_CLOSING_CHANNEL, t)
	expectSmtpStatusText(commandResponse, "Single message per session only. Closing session.", t)
}

func TestAttemptDataBeforeRcpt(t *testing.T) {

	s := NewCommandProcessor()

	mockChannel := make(chan *SmtpServerChannelMessage)
	smtpServerInstance := SmtpServer_impl{config: common.CreateDefaultAppConfig(), channel: mockChannel}
	smtpServerInstance.config.GetSmtpServerConfig().SMTP_OPTION_SINGLE_MESSAGE_PER_SESSION = false

	clientConn, _ := getTestConnection()

	mockSession := NewSmtpSession(smtpServerInstance, clientConn, mockChannel)

	mockSession.beginSession()

	commandResponse := s.HandleCommand(mockSession, "mail from: a@a.com")

	expectSmtpSessionAction(commandResponse, COMMANDACTION_CONTINUE, t)
	expectSmtpStatusCode(commandResponse, SMTP_COMMAND_STATUS_BAD_COMMAND_SEQUENCE, t)
	expectSmtpStatusText(commandResponse, "Send EHLO/HELO first!", t)
}

func TestAttemptDataBeforeMailFrom(t *testing.T) {

	s := NewCommandProcessor()

	mockChannel := make(chan *SmtpServerChannelMessage)
	smtpServerInstance := SmtpServer_impl{config: common.CreateDefaultAppConfig(), channel: mockChannel}
	smtpServerInstance.config.GetSmtpServerConfig().SMTP_OPTION_SINGLE_MESSAGE_PER_SESSION = false

	clientConn, _ := getTestConnection()

	mockSession := NewSmtpSession(smtpServerInstance, clientConn, mockChannel)

	mockSession.beginSession()

	commandResponse := s.HandleCommand(mockSession, "mail from: a@a.com")

	expectSmtpSessionAction(commandResponse, COMMANDACTION_CONTINUE, t)
	expectSmtpStatusCode(commandResponse, SMTP_COMMAND_STATUS_BAD_COMMAND_SEQUENCE, t)
	expectSmtpStatusText(commandResponse, "Send EHLO/HELO first!", t)
}

func TestDisallowBareLineFeedAfterData(t *testing.T) {

	s := NewCommandProcessor()

	mockChannel := make(chan SmtpServerChannelMessage)
	smtpServerInstance := SmtpServer_impl{config: common.CreateDefaultAppConfig(), channel: mockChannel}
	smtpServerInstance.config.GetSmtpServerConfig().SMTP_OPTION_SINGLE_MESSAGE_PER_SESSION = false

	clientConn, _ := getTestConnection()

	mockSession := NewSmtpSession(smtpServerInstance, clientConn, mockChannel)

	mockSession.beginSession()
	s.HandleCommand(mockSession, "helo 127.0.0.1")

	s.HandleCommand(mockSession, "mail from: a@a.com")
	s.HandleCommand(mockSession, "rcpt to: a@a.com")
	s.HandleCommand(mockSession, "data")
	commandResponse := s.HandleCommand(mockSession, "To: a@a.com\nFrom: a@a.com\nSubject: This is a test\n\n.\n")

	expectSmtpSessionAction(commandResponse, COMMANDACTION_CONTINUE, t)
	expectSmtpStatusCode(commandResponse, SMTP_COMMAND_STATUS_REQUESTED_ACTION_ABORTED, t)
	expectSmtpStatusText(commandResponse, "Non CRLF line-feed submit not supported", t)

}
