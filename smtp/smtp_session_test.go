package smtp

import (
	"github.com/andrewjc/milhaux/common"
	"github.com/petergtz/pegomock"
	"testing"
)

func TestConnectingToServerPresentsWelcomeMessage(t *testing.T) {
	pegomock.RegisterMockTestingT(t)

	mockChannel := make(chan SmtpServerChannelMessage)
	smtpServerInstance := SmtpServer_impl{config: common.CreateDefaultAppConfig(), channel: mockChannel}

	clientConn, _ := getTestConnection()

	mockSession := NewSmtpSession(smtpServerInstance, clientConn, mockChannel)

	if err := mockSession.beginSession(); err != nil {
		t.Errorf("beginSession returned an error: %v", err.Error())
	}
}
