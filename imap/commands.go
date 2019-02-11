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
