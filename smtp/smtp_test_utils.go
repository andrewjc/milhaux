package smtp

import (
	"testing"
)

func expectSmtpSessionAction(response CommandResponse, expectedSmtpSessionAction CommandAction, t *testing.T) {
	if response.action != expectedSmtpSessionAction {
		t.Errorf("Expected %v - got %v", commandActionToString(expectedSmtpSessionAction), response.action.String())
	}
}

func expectSmtpStatusCode(response CommandResponse, expectedSmtpCommandStatusCode SmtpCommandStatus, t *testing.T) {
	if response.commandStatus != expectedSmtpCommandStatusCode {
		t.Errorf("Expected %v - got %v", expectedSmtpCommandStatusCode, response.commandStatus)
	}
}

func expectSmtpStatusText(response CommandResponse, expectedSmtpCommandStatusText string, t *testing.T) {
	if response.commandResponseText != expectedSmtpCommandStatusText {
		t.Errorf("Expected different command response text '%v' - got '%v'", expectedSmtpCommandStatusText, response.commandResponseText)
	}
}

func commandActionToString(action CommandAction) string {
	s := CommandResponse{}
	s.action = action
	return s.action.String()
}
