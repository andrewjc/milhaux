package smtp

import (
	"strings"
)

type CommandAction uint16
type CommandStatus uint16

const (
	COMMANDACTION_CONTINUE = iota
	COMMANDACTION_EXIT
	COMMANDACTION_NONE

	END_DATA_COMMAND_SEQUENCE = ".\r\n"
)

type CommandResponse struct {
	action              CommandAction
	commandStatus       CommandStatus
	commandResponseText string
}

func (s *SmtpServer_impl) commandRequestHandler(smtpSession *SmtpSession, commandLine string) *CommandResponse {

	command := getCommandArgPair(commandLine)

	if command.commandStr == SMTP_COMMAND_DONE {
		return s.smtpCommandDone(smtpSession, commandLine)
	}

	if strings.TrimSpace(string(command.commandStr)) == "EXIT" {
		return &CommandResponse{COMMANDACTION_EXIT, SMTP_COMMAND_STATUS_SERVICE_READY, "BYE"}
	}

	switch {
	case smtpSession.SmtpState == SMTP_SERVER_STATE_ESTABLISH:
		return s.smtpCommandEstablish(smtpSession, command)
	case smtpSession.StateData[SESSION_DATA_KEY_CLIENT_ID] != nil:

		if smtpSession.SmtpState == SMTP_SERVER_STATE_DATA {
			return s.smtpCommandBufferData(smtpSession, commandLine)
		} else {
			switch {
			case command.commandStr == SMTP_COMMAND_MAIL:
				return s.smtpCommandMail(smtpSession, command)
			case command.commandStr == SMTP_COMMAND_RCPT:
				return s.smtpCommandMail(smtpSession, command)
			case command.commandStr == SMTP_COMMAND_DATA:
				return s.smtpCommandReceiveData(smtpSession, command)
			}
		}

	}

	return &CommandResponse{COMMANDACTION_EXIT, SMTP_COMMAND_STATUS_COMMAND_NOT_IMPLEMENTED, "Unknown command"}
}

type commandArgPair struct {
	commandStr string
	args       string
}

/* Utils */
func getCommandArgPair(rawString string) commandArgPair {

	temp := strings.TrimSpace(rawString)

	if strings.ContainsAny(temp, " ") {

		commandStr := temp[0:strings.Index(temp, " ")]
		argStr := temp[len(commandStr):len(temp)]
		return commandArgPair{strings.TrimSpace(commandStr), strings.TrimSpace(argStr)}
	} else {
		return commandArgPair{temp, ""}
	}
}
