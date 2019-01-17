package smtp

import (
	"net"
	"testing"
)

func expectSmtpSessionAction(response *CommandResponse, expectedSmtpSessionAction CommandAction, t *testing.T) {
	if response.action != expectedSmtpSessionAction {
		t.Errorf("Expected command %v. Got: %v %v", commandActionToString(expectedSmtpSessionAction), response.actionToString(), response.commandStatus)
	}
}

func expectSmtpStatusCode(response *CommandResponse, expectedSmtpCommandStatusCode CommandStatus, t *testing.T) {
	if response.commandStatus != expectedSmtpCommandStatusCode {
		t.Errorf("Expected %v ok", expectedSmtpCommandStatusCode)
	}
}

func expectSmtpStatusText(response *CommandResponse, expectedSmtpCommandStatusText string, t *testing.T) {
	if response.commandResponseText != expectedSmtpCommandStatusText {
		t.Errorf("Expected different command response text '%v' - got '%v'", expectedSmtpCommandStatusText, response.commandResponseText)
	}
}

func commandActionToString(action CommandAction) string {
	s := CommandResponse{}
	s.action = action
	return s.actionToString()
}

func getTestConnection() (client, server net.Conn) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")

	var serverConn net.Conn
	go func() {
		defer ln.Close()
		server, _ = ln.Accept()
	}()

	clientConn, _ := net.Dial("tcp", ln.Addr().String())

	return clientConn, serverConn
}
