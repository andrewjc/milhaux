package backend

import (
	"errors"
	"fmt"
	"github.com/andrewjc/milhaux/common"
)

type BackendCommandProcessor struct{}

func (m *MailStoreBackend) handleConnectorMessage(message common.StorageMessage) (common.StorageMessage, error) {
	// route the message
	switch {
	case message.MsgType() == common.StorageMessageTypeQuery:
		return m.handleQueryMessage(message)
	}

	return nil, errors.New(fmt.Sprintf("Unhandled backend connector message: %s", message.ToJson()))
}

func (m *MailStoreBackend) handleQueryMessage(message common.StorageMessage) (common.StorageMessage, error) {

	switch {
	case message.MsgCommand() == common.StorageMessageCommandIsValid:
		isValid := m.commandProcessor.IsValidCommand(message.MsgData())
		responseMessageBody := "false"
		if isValid {
			responseMessageBody = "true"
		}
		return (&common.StorageMessageBuilder{}).ResponseMessage(string(responseMessageBody)).Build(), nil
	}

	return nil, errors.New(fmt.Sprintf("Unhandled backend connector message: %s", message.ToJson()))
}
