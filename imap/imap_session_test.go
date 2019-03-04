package imap

import (
	"github.com/andrewjc/milhaux/common"
	apptesting "github.com/andrewjc/milhaux/testing"
	"github.com/petergtz/pegomock"
	"testing"
)

func TestConnectAuthToImapServer(t *testing.T) {
	pegomock.RegisterMockTestingT(t)

	s := NewCommandProcessor()

	imapServerInstance := Imap4Server_Impl{config: common.CreateDefaultAppConfig(), commandProcessor:NewCommandProcessor()}

	clientConn, _ := apptesting.GetTestConnection()

	session := NewImapSession(imapServerInstance, clientConn)

	session.beginSession()
	resp := s.HandleCommand(session, "AJC123 LOGIN andrew.cranston@gmail.com password123")

	expectSessionAction(resp, COMMANDACTION_CONTINUE, t)
	expectStatusCode(resp, IMAP_COMMAND_STATUS_OK, t)
	expectResponseLines(resp, []string{"* CAPABILITY IMAP4rev1 AUTH=PLAIN", "AJC123 OK andrew authenticated (Success)"}, t)
}

func TestConnectAuthToImapServerIncorrectPassword(t *testing.T) {
	pegomock.RegisterMockTestingT(t)

	s := NewCommandProcessor()

	imapServerInstance := Imap4Server_Impl{config: common.CreateDefaultAppConfig(), commandProcessor:NewCommandProcessor()}

	clientConn, _ := apptesting.GetTestConnection()

	session := NewImapSession(imapServerInstance, clientConn)

	session.beginSession()
	resp := s.HandleCommand(session, "AJC123 LOGIN andrew.cranston@gmail.com password123fff")

	expectSessionAction(resp, COMMANDACTION_CONTINUE, t)
	expectStatusCode(resp, IMAP_COMMAND_STATUS_OK, t)
	expectResponseLines(resp, []string{"* CAPABILITY IMAP4rev1 AUTH=PLAIN", "AJC123 OK andrew authenticated (Success)"}, t)
}
