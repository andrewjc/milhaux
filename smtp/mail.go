package smtp

import (
	"fmt"
	"github.com/andrewjc/milhaux/common"
	"strings"
)

func (s *SmtpServer_impl) smtpCommandMail(session *SmtpSession, commandArg commandArgPair) *CommandResponse {

	switch {
	case strings.HasPrefix(commandArg.args, "FROM:"):
		session.StateData[SESSION_DATA_KEY_MAIL_FROM] = string(commandArg.args[len("FROM:")])
		session.SmtpState = SMTP_SERVER_STATE_MAIL

		return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_MAIL_ACTION_OK, "OK"}
	case strings.HasPrefix(commandArg.args, "TO:"):
		session.StateData[SESSION_DATA_KEY_MAIL_TO] = string(commandArg.args[len("TO:")])
		session.SmtpState = SMTP_SERVER_STATE_RCPT
		return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_MAIL_ACTION_OK, "OK"}
	default:
		return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_COMMAND_NOT_IMPLEMENTED, fmt.Sprintf("Unknown command: %s %s", commandArg.commandStr, commandArg.args)}
	}
}

func (s *SmtpServer_impl) smtpCommandReceiveData(session *SmtpSession, commandArg commandArgPair) *CommandResponse {

	session.SmtpState = SMTP_SERVER_STATE_DATA
	return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_START_MAIL_INPUT, "Start mail input; end with <CRLF>.<CRLF>"}

}

func (s *SmtpServer_impl) smtpCommandBufferData(smtpSession *SmtpSession, commandLine string) *CommandResponse {
	if commandLine == END_DATA_COMMAND_SEQUENCE {
		newMessage := &common.MailMessage{
			To:   smtpSession.StateData[SESSION_DATA_KEY_MAIL_FROM].(string),
			From: smtpSession.StateData[SESSION_DATA_KEY_MAIL_TO].(string),
			Data: smtpSession.ReceiveDataBuffer.String(),
		}
		s.onSubmitMail(newMessage)
		smtpSession.SmtpState = SMTP_SERVER_STATE_DONE
		return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_MAIL_ACTION_OK, fmt.Sprintf("OK. Queued for Delivery. Queue id: %s", newMessage.QueueId)}
	} else {
		smtpSession.ReceiveDataBuffer.WriteString(commandLine)
		return &CommandResponse{COMMANDACTION_NONE, SMTP_COMMAND_STATUS_NONE, ""}
	}
}
