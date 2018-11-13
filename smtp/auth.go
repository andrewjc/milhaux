package smtp

import (
	"fmt"
)

func (s *SmtpServer_impl) smtpCommandEstablish(smtpSession *SmtpSession, command commandArgPair) *CommandResponse {
	switch command.commandStr {
	case "HELO": //hello command
		return s.smtpCommandHelo(smtpSession, command)
	case "EHLO": //enhanced hello
		return s.smtpCommandEhlo(smtpSession, command)
	case "STARTTLS": // start authenticated tls session
		return s.smtpCommandStartTls(smtpSession, command)
	default: // wtf is this?
		return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_COMMAND_NOT_IMPLEMENTED, fmt.Sprintf("Unknown command: %s %s", command.commandStr, command.args)}
	}
}

func (s *SmtpServer_impl) smtpCommandHelo(session *SmtpSession, commandArg commandArgPair) *CommandResponse {

	session.SmtpState = SMTP_SERVER_STATE_GOTHELO
	session.StateData[SESSION_DATA_KEY_CLIENT_ID] = commandArg.args

	return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_MAIL_ACTION_OK, fmt.Sprintf("%s Hello", commandArg.args)}
}

func (s *SmtpServer_impl) smtpCommandEhlo(session *SmtpSession, commandArg commandArgPair) *CommandResponse {

	session.SmtpState = SMTP_SERVER_STATE_GOTHELO
	session.StateData[SESSION_DATA_KEY_CLIENT_ID] = commandArg.args

	return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_MAIL_ACTION_OK, fmt.Sprintf("%s Hello", commandArg.args)}
}

func (s *SmtpServer_impl) smtpCommandStartTls(session *SmtpSession, commandArg commandArgPair) *CommandResponse {

	session.SmtpState = SMTP_SERVER_STATE_GOTHELO
	session.StateData[SESSION_DATA_KEY_CLIENT_ID] = commandArg.args

	return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_MAIL_ACTION_OK, fmt.Sprintf("%s Hello", commandArg.args)}
}
