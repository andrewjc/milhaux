package imap

import (
	"fmt"
	"github.com/andrewjc/milhaux/common"
	"strings"
)

func (s *ImapCommandProcessor) imapCommandCapability(session *ImapSession, commandArg commandArgPair) CommandResponse {
	commandOkLine := fmt.Sprintf("%s %s %s %s", commandArg.sequenceIdent, IMAP_COMMAND_STATUS_OK, "That's all!", session.sessionToken)

	response := CommandResponse{}
	response.action = COMMANDACTION_CONTINUE
	response.commandStatus = IMAP_COMMAND_STATUS_OK
	response.commandResponseLines = []string{strings.TrimSpace(s.getCapabilityString()), commandOkLine}

	return response
}

func (s *ImapCommandProcessor) imapCommandLogout(session *ImapSession, commandArg commandArgPair) CommandResponse {
	cLineBuilder := strings.Builder{}
	cLineBuilder.WriteString("* BYE")

	cLineBuilder.WriteString(" LOGOUT user requested ")

	commandOkLine := fmt.Sprintf("%s %s %s %s", commandArg.sequenceIdent, IMAP_COMMAND_STATUS_OK, "Come again soon!", session.sessionToken)

	response := CommandResponse{}
	response.action = COMMANDACTION_EXIT
	response.commandStatus = IMAP_COMMAND_STATUS_OK
	response.commandResponseLines = []string{cLineBuilder.String(), commandOkLine}

	return response
}

func (s *ImapCommandProcessor) imapCommandNoop(session *ImapSession, commandArg commandArgPair) CommandResponse {
	commandOkLine := fmt.Sprintf("%s %s %s %s", commandArg.sequenceIdent, IMAP_COMMAND_STATUS_OK, "Nothing Accomplished.", session.sessionToken)

	response := CommandResponse{}
	response.action = COMMANDACTION_CONTINUE
	response.commandStatus = IMAP_COMMAND_STATUS_OK
	response.commandResponseLines = []string{commandOkLine}

	return response
}

func (s *ImapCommandProcessor) imapCommandLogin(session *ImapSession, commandArg commandArgPair) CommandResponse {
	fields := common.ParseQuoteAwareFields(commandArg.args, 2)

	username := fields[0]
	//password := fields[1]

	commandOkLine := fmt.Sprintf("%s %s %s %s", commandArg.sequenceIdent, IMAP_COMMAND_STATUS_OK, username, "authenticated (Success)")

	response := CommandResponse{}
	response.action = COMMANDACTION_CONTINUE
	response.commandStatus = IMAP_COMMAND_STATUS_OK
	response.commandResponseLines = []string{strings.TrimSpace(s.getUserCapabilityString()), commandOkLine}

	session.sessionState = IMAP_SESSION_STATE_AUTHOK

	return response
}

func (s *ImapCommandProcessor) getCapabilityString() string {
	cLineBuilder := strings.Builder{}
	cLineBuilder.WriteString("* CAPABILITY IMAP4rev1 AUTH=PLAIN")

	return cLineBuilder.String()
}

// TODO: This needs to return capabilities that are granted to the user logging in...
func (s *ImapCommandProcessor) getUserCapabilityString() string {
	cLineBuilder := strings.Builder{}
	cLineBuilder.WriteString("* CAPABILITY IMAP4rev1 AUTH=PLAIN")

	return cLineBuilder.String()
}

func (s *ImapCommandProcessor) isImapStorageDriverCommand(imapSession *ImapSession, commandPair commandArgPair) bool {

	storageConnector := imapSession.imapServerInstance.storageConnector

	message := (&common.StorageMessageBuilder{}).IsValidCommandMessage(commandPair.commandStr.String()).Build()

	responseMessage, _ := storageConnector.PerformSendReceive(message)

	if responseMessage.MsgType() == common.StorageMessageTypeResponse {
		return responseMessage.MsgData() == "true"
	}

	return false
}

// This is a command that needs to be delegated to the backend.
// The storage connector will pass the message for us.
func (s *ImapCommandProcessor) imapStorageDriverCommandHandler(imapSession *ImapSession, commandPair commandArgPair) CommandResponse {
	storageConnector := imapSession.imapServerInstance.storageConnector

	message := (&common.StorageMessageBuilder{}).IsValidCommandMessage(commandPair.commandStr.String()).Build()

	responseMessage, _ := storageConnector.PerformSendReceive(message)

	if responseMessage.MsgType() == common.StorageMessageTypeResponse {
		// TODO: The response message from the backend should indicate status and status text...

		// for example, if we get RENAME "Inbox/Demo12" "Inbox/Demo1" from the client
		// we send a rename message to the backend
		// if we get back a success message, then we set action to continue

		// only special cases need to be considered for when action = exit

		status := IMAP_COMMAND_STATUS_OK
		action := COMMANDACTION_CONTINUE
		responseLines := []string{"TODO!!!"}
		return CommandResponse{commandStatus: status, action: action, commandResponseLines: responseLines}
	}

	return CommandResponse{commandStatus: IMAP_COMMAND_STATUS_COMMAND_NOT_IMPLEMENTED, action: COMMANDACTION_CONTINUE, commandResponseLines: []string{fmt.Sprintf("Imap storage driver received an unhandled connector response: %s", responseMessage)}}
}
