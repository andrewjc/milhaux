package smtp

import (
	"strings"
)

const (
	END_DATA_COMMAND_SEQUENCE         = ".\r\n"
	END_DATA_INVALID_COMMAND_SEQUENCE = "\n\n.\n"
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

type SmtpCommandProcessor struct {
}

func NewCommandProcessor() SmtpCommandProcessor {
	processor := SmtpCommandProcessor{}

	return processor
}

type CommandResponse struct {
	action              CommandAction
	commandStatus       SmtpCommandStatus
	commandResponseText string
}

func (s *SmtpCommandProcessor) HandleCommand(smtpSession *SmtpSession, commandLine string) CommandResponse {

	command := getCommandArgPair(commandLine)

	if strings.TrimSpace(string(command.commandStr)) == "EXIT" {
		return CommandResponse{COMMANDACTION_EXIT, SMTP_COMMAND_STATUS_SERVICE_READY, "BYE"}
	}

	switch {
	case smtpSession.smtpState == SMTP_SESSION_STATE_PREAUTH:
		return s.smtpCommandEstablish(smtpSession, command)
	case smtpSession.stateData[SESSION_DATA_KEY_CLIENT_ID] != nil:

		if smtpSession.smtpState == SMTP_SESSION_STATE_DATA {
			return s.smtpCommandBufferData(smtpSession, commandLine)
		} else {

			// single or multi message per session policy
			if smtpSession.smtpServerInstance.config.GetSmtpServerConfig().SMTP_OPTION_SINGLE_MESSAGE_PER_SESSION == true {
				if smtpSession.smtpState == SMTP_SESSION_STATE_SUBMIT {
					// A session has already submitted a message
					return CommandResponse{COMMANDACTION_EXIT, SMTP_COMMAND_STATUS_SERVICE_CLOSING_CHANNEL, "Single message per session only. Closing session."}
				}
			}

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

	return CommandResponse{COMMANDACTION_EXIT, SMTP_COMMAND_STATUS_COMMAND_NOT_IMPLEMENTED, "Unknown command"}
}

type commandArgPair struct {
	commandStr SmtpCommandVerb
	args       string
}

/* Utils */
func getCommandArgPair(rawString string) commandArgPair {

	temp := strings.TrimSpace(rawString)

	if strings.ContainsAny(temp, " ") {

		commandStr := temp[0:strings.Index(temp, " ")]
		argStr := temp[len(commandStr):len(temp)]
		return commandArgPair{SmtpCommandVerb(strings.ToUpper(strings.TrimSpace(commandStr))), strings.TrimSpace(argStr)}
	} else {
		return commandArgPair{SmtpCommandVerb(strings.ToUpper(temp)), ""}
	}
}
