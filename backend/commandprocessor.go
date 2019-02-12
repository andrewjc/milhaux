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
		isValid := m.commandProcessor.IsValidCommand(message.MsgData().(string))
		return (&common.StorageMessageBuilder{}).SuccessResponseMessage(common.Bool2String(isValid)).Build(), nil
	case message.MsgCommand() == common.StorageMessageCommandList:
		return m.handleListCommand(message)
	case message.MsgCommand() == common.StorageMessageCommandCreate:
		return m.handleCreateCommand(message)
	}

	return nil, errors.New(fmt.Sprintf("Unhandled backend connector message: %s", message.ToJson()))
}

func (m *MailStoreBackend) handleListCommand(message common.StorageMessage) (common.StorageMessage, error) {
	args := common.ParseQuoteAwareFields(message.MsgData().(string), 2)

	referenceName := common.StripQuotes(args[0])
	mailboxName := common.StripQuotes(args[1])

	if referenceName == "" {
		// mailbox name is a SELECT
		if mailboxName == "" {
			// return the delimiter
			lines := []string{
				"* LIST (\\Noselect) \"/\" \"\"",
			}
			return (&common.StorageMessageBuilder{}).SuccessResponseMessage(lines).Build(), nil
		}

		if mailboxName == "*" {
			// return all names
			var listSlice []string
			for _, x := range m.storageProvider.Mailboxes() {
				listSlice = append(listSlice, "* LIST () \"/\" \""+x.Name()+"\"")
			}
			return (&common.StorageMessageBuilder{}).SuccessResponseMessage(listSlice).Build(), nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Bad list command: %s", message.ToJson()))
}

func (m *MailStoreBackend) handleCreateCommand(message common.StorageMessage) (common.StorageMessage, error) {
	args := common.ParseQuoteAwareFields(message.MsgData().(string), 1)

	mailboxName := common.StripQuotes(args[0])

	if mailboxName == "" {
		return nil, errors.New(fmt.Sprintf("Bad create command: %s", message.ToJson()))
	} else {
		// return all names
		return (&common.StorageMessageBuilder{}).SuccessResponseMessage(nil).Build(), nil
	}

}
