package imap

import (
	"errors"
	"fmt"
	"strings"
)

type CommandAction uint16

const (
	COMMANDACTION_CONTINUE = CommandAction(1)
	COMMANDACTION_EXIT     = CommandAction(2)
	COMMANDACTION_NONE     = CommandAction(3)
)

func (s CommandAction) String() string {
	switch s {
	case COMMANDACTION_CONTINUE:
		return "CONTINUE"
	case COMMANDACTION_EXIT:
		return "EXIT"
	case COMMANDACTION_NONE:
		return "NONE"
	default:
		return "Unknown"
	}
}

type ImapCommandProcessor struct {
}

func NewCommandProcessor() ImapCommandProcessor {
	processor := ImapCommandProcessor{}

	return processor
}

type CommandResponse struct {
	action               CommandAction
	commandStatus        ImapCommandStatus
	commandResponseLines []string
}

func (s *ImapCommandProcessor) HandleCommand(imapSession *ImapSession, commandLine string) CommandResponse {

	command, err := getCommandArgPair(commandLine)

	if err != nil {
		return CommandResponse{COMMANDACTION_EXIT, IMAP_COMMAND_STATUS_COMMAND_NOT_IMPLEMENTED, []string{fmt.Sprintf("Invalid command sequence: %s", err.Error())}}
	}

	switch {
	case command.commandStr == IMAP_COMMAND_CAPABILITY:
		return s.imapCommandCapability(imapSession, command)
	case command.commandStr == IMAP_COMMAND_LOGOUT:
		return s.imapCommandLogout(imapSession, command)
	case command.commandStr == IMAP_COMMAND_NOOP:
		return s.imapCommandNoop(imapSession, command)
	}

	// Commands that are only allowed before auth is performed
	if imapSession.smtpState == IMAP_SESSION_STATE_PREAUTH {
		switch {
		case command.commandStr == IMAP_COMMAND_LOGIN:
			return s.imapCommandLogin(imapSession, command)
		}
	}

	return CommandResponse{COMMANDACTION_EXIT, IMAP_COMMAND_STATUS_COMMAND_NOT_IMPLEMENTED, []string{"Unknown command"}}
}

type commandArgPair struct {
	sequenceIdent string
	commandStr    ImapCommandVerb
	args          string
}

/* Utils */
func getCommandArgPair(rawString string) (commandArgPair, error) {

	pair := commandArgPair{}

	temp := strings.TrimSpace(rawString)

	// the rawstring needs to have atleast 1 space, otherwise it's just a sequence id with no commands or args
	if strings.ContainsAny(temp, " ") {
		fields := strings.Fields(temp)
		pair.sequenceIdent = fields[0]

		commandStr := fields[1]
		pair.commandStr = ImapCommandVerb(commandStr)
		if len(fields) > 2 {
			argStr := temp[(len(pair.sequenceIdent) + len(commandStr) + 2):len(temp)]
			pair.args = argStr
		}

		return pair, nil

	} else {
		return pair, errors.New("invalid command sequence")
	}
}
