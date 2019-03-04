package imap

import (
	"testing"
)

func expectSessionAction(response CommandResponse, expectedSmtpSessionAction CommandAction, t *testing.T) {
	if response.action != expectedSmtpSessionAction {
		t.Errorf("Expected %v - got %v", commandActionToString(expectedSmtpSessionAction), response.action.String())
	}
}

func expectStatusCode(response CommandResponse, expectedSmtpCommandStatusCode ImapCommandStatus, t *testing.T) {
	if response.commandStatus != expectedSmtpCommandStatusCode {
		t.Errorf("Expected %v - got %v", expectedSmtpCommandStatusCode, response.commandStatus)
	}
}

func expectResponseLines(response CommandResponse, expectedResponseLines []string, t *testing.T) {
	for x, y := range expectedResponseLines {
		if response.commandResponseLines[x] != y {
			t.Errorf("Expected different command response text '%v' - got '%v'", expectedResponseLines, response.commandResponseLines)
		}
	}
}

func commandActionToString(action CommandAction) string {
	s := CommandResponse{}
	s.action = action
	return s.action.String()
}
